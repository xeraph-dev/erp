**another question
I want to allow multiple Content-Type, so if the user request Content-Type is application/json, the response will be application/json, same for text/xml (right now only those are implemented, but more later)
what about that decision?**

---

## 2. Timeout middleware

Wrap the request context with a deadline and let it propagate — this is the main win, since `pgx` respects context cancellation, so a timed-out request will actually abort its in-flight DB query rather than leaving it running.

```go
// internal/middlewares/timeout_middleware.go
package middlewares

import (
	"context"
	"net/http"
	"time"
)

func Timeout(d time.Duration) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx, cancel := context.WithTimeout(r.Context(), d)
			defer cancel()

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
```

This alone gets the cancellation propagation, but doesn't stop a handler that ignores `ctx.Done()` from writing a response after the deadline. For a hard boundary — actually returning `503` to the client when the deadline hits, and preventing writes afterward — wrap with the stdlib's `http.TimeoutHandler` instead of hand-rolling it:

```go
func Timeout(d time.Duration) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.TimeoutHandler(next, d, `{"error":"request timed out"}`)
	}
}
```

`http.TimeoutHandler` runs your handler in a goroutine, races it against a timer, and — if the timer wins — writes `503` with the given body and blocks further writes from the (now-abandoned) handler goroutine. Your inner handler still needs to respect `ctx.Done()` to actually stop working (e.g., a long-running query), so the two approaches aren't mutually exclusive — `TimeoutHandler` for the client-facing contract, `context.WithTimeout` propagation for cooperative cancellation of the actual work. I'd use `TimeoutHandler`, since it subsumes the manual version and gives you the client response for free.

Register per-route rather than globally if some endpoints legitimately need longer (report generation, bulk imports):

```go
server.Use(
	middlewares.Recoverer,
	middlewares.Logger(logger),
	middlewares.RequestID,
	middlewares.HTTPLogger,
	middlewares.Timeout(10 * time.Second),
)
```

or scope it inside a `server.Group(...)` for just the slow endpoints.

## 3. Rate limiting by IP

Since Redis is already in your stack, I'd use it rather than an in-memory limiter — an in-memory map only rate-limits per-instance, so the moment you run more than one replica behind nginx, each instance tracks its own counts and the effective limit multiplies by instance count. If you're confident you'll stay single-instance for a long while, in-memory is simpler; otherwise go straight to Redis to avoid a painful migration later.

**In-memory version** (`golang.org/x/time/rate`, single instance only):

```go
// internal/middlewares/rate_limit_middleware.go
package middlewares

import (
	"net"
	"net/http"
	"sync"

	"golang.org/x/time/rate"
)

type ipLimiter struct {
	mu       sync.Mutex
	limiters map[string]*rate.Limiter
	r        rate.Limit
	burst    int
}

func newIPLimiter(r rate.Limit, burst int) *ipLimiter {
	return &ipLimiter{limiters: make(map[string]*rate.Limiter), r: r, burst: burst}
}

func (l *ipLimiter) get(ip string) *rate.Limiter {
	l.mu.Lock()
	defer l.mu.Unlock()

	lim, ok := l.limiters[ip]
	if !ok {
		lim = rate.NewLimiter(l.r, l.burst)
		l.limiters[ip] = lim
	}
	return lim
}

func RateLimit(requestsPerSecond float64, burst int) func(next http.Handler) http.Handler {
	limiter := newIPLimiter(rate.Limit(requestsPerSecond), burst)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip, _, err := net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				ip = r.RemoteAddr
			}

			if !limiter.get(ip).Allow() {
				w.WriteHeader(http.StatusTooManyRequests)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
```

Note this map grows unbounded — every distinct IP that ever hits your server gets an entry that's never evicted. Fine for a low-traffic internal ERP tool, but add a periodic sweep (or an LRU cache like `hashicorp/golang-lru`) if you expect real internet traffic.

**Redis version** (works correctly across multiple instances — recommended given your stack):

```go
// internal/middlewares/rate_limit_middleware.go
package middlewares

import (
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
)

func RateLimit(rdb *redis.Client, limit int64, window time.Duration) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			logger := GetLogger(ctx)

			ip, _, err := net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				ip = r.RemoteAddr
			}
			key := fmt.Sprintf("ratelimit:%s", ip)

			count, err := rdb.Incr(ctx, key).Result()
			if err != nil {
				logger.ErrorContext(ctx, "rate limit check failed", "error", err)
				next.ServeHTTP(w, r) // fail open: don't block traffic on Redis being down
				return
			}
			if count == 1 {
				rdb.Expire(ctx, key, window)
			}

			if count > limit {
				w.WriteHeader(http.StatusTooManyRequests)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
```

`INCR` + `EXPIRE` gives you a fixed-window counter — simple and correct for most abuse-prevention needs. It's not perfectly smooth at window boundaries (a client could burst near the edge of two windows), but that's rarely worth the complexity of a sliding-window or token-bucket algorithm in Redis unless you have a specific reason to need it.

**One important caveat on `r.RemoteAddr`:** if nginx sits in front of your Go service (per your stack), `r.RemoteAddr` will be nginx's address, not the real client's — you need `X-Forwarded-For` or `X-Real-IP`, set by nginx, and you must only trust that header if nginx is guaranteed to be the sole entry point (otherwise a client can spoof it and dodge the limit):

```go
ip := r.Header.Get("X-Real-IP")
if ip == "" {
	ip, _, _ = net.SplitHostPort(r.RemoteAddr)
}
```

Configure nginx to set `X-Real-IP $remote_addr;` in its proxy config so this is reliable.

## Putting it together

```go
server.Use(
	middlewares.Recoverer,
	middlewares.Logger(logger),
	middlewares.RequestID,
	middlewares.HTTPLogger,
	middlewares.RateLimit(redisClient, 100, time.Minute),
	middlewares.Timeout(10 * time.Second),
)
```

Order matters: `Recoverer` outermost (as established earlier), rate limiting before timeout so a limited request gets rejected immediately without consuming a timeout slot, and both before the request reaches your handlers.

---

**my idea of roles level is that a user only can assign a role of lower level, and new users are user role
so, following that constraint, the system only can have one admin user and the user user is the lowest level**

---

Your rule is logically sound, but it has a consequence worth stating explicitly before you build on it: **strict "lower than my own level" means no one — not even the admin — can ever create a second admin.** To grant a role at level _L_, the assigner needs level > _L_. For _L_ = admin's level (32767, the max), no level exists above that, so the constraint mathematically forbids granting admin to anyone, ever, through the normal flow. Your "system can only have one admin" isn't a rule you need to separately enforce — it falls out of the level comparison for free. Good catch on your part, but it creates a real operational risk that needs a deliberate answer, not an accidental one.

## The single point of failure this creates

If the one admin account is lost — credentials forgotten, employee leaves, account compromised and locked — **there is no in-app path to recover admin access.** No user has a level high enough to grant it. You need one of these, chosen deliberately:

1. **Accept it, and build an out-of-band recovery path.** A CLI tool or direct SQL script (run with your `migrator` role, outside the API) that inserts a `roles_users` row bypassing the application-layer check entirely. This is consistent with your "provisioning scripts fail loud, run once" philosophy — treat admin recovery the same way: an explicit, logged, deliberate operation outside normal request flow, not a hidden API backdoor.
2. **Loosen the rule for the admin level specifically** — allow assigning a role _equal to_ the assigner's own level only when that level is the max (admin can create peer admins), while every other level keeps the strict `<` rule. This gives you resilience without breaking the "can't create peers below admin" property you actually want.

I'd lean toward (1): it keeps the invariant clean and auditable, and matches your existing pattern of treating dangerous operations as deliberate, out-of-band actions rather than API-reachable ones.

## Effective level with multi-role users

`roles_users` is many-to-many, so a user can hold multiple roles simultaneously. "The assigner's level" needs a precise definition — I'd use **`MAX(role_level)` across all of the assigner's active (non-soft-deleted) roles**, not a single role field:

```sql
SELECT MAX(r.role_level)
FROM erp.roles ru_roles  -- the view, so soft-deletes are already filtered
JOIN erp.roles_users ru ON ru.role_id = ru_roles.id
JOIN erp.roles r ON r.id = ru.role_id
WHERE ru.user_id = $1;
```

## Apply the same rule to revocation, not just assignment

If Alice (level 100) can't _grant_ a level-500 role, she also shouldn't be able to _revoke_ one from someone else — otherwise a mid-level user could strip a higher-level user's role even though they could never have granted it. Mirror the check symmetrically in both directions.

## Where to enforce it

**Service layer, using the authenticated user's resolved level from the session — never from client input.** This is the natural home given your existing `RoleService`/`UserService` split:

```go
func (service RoleServiceImpl) AssignRole(ctx context.Context, assignerID, targetUserID, roleID uuid.UUID) (err error) {
	assignerLevel, err := service.Repos.Role.MaxLevelForUser(ctx, service.DB, assignerID)
	if err != nil {
		return err
	}

	targetRole, err := service.Repos.Role.GetByID(ctx, service.DB, roleID)
	if err != nil {
		return err
	}

	if targetRole.Level >= assignerLevel {
		return ErrInsufficientPrivilege
	}

	// proceed with withTx-wrapped assignment
}
```

`assignerID` should come from the authenticated context (your auth middleware already propagates user ID) — never trust a client-supplied "my level is X" claim.

**Add a DB trigger as defense-in-depth**, consistent with your existing audit/soft-delete trigger pattern — a `BEFORE INSERT` trigger on `roles_users` that re-checks the same invariant using `app.current_user_id` (already set in your transaction helper). This protects you if a bug ever lets a bad assignment slip past the service layer, or if something writes to `roles_users` outside the service (a future admin script, a migration, etc.).

## New user default role

Insert the "user" role assignment **inside the same `withTx` transaction as user creation** in `Register`, not as a separate follow-up call — if the role insert fails, the whole registration should roll back rather than leaving a roleless user. Since `roles_users.created_by_id` defaults to `erp.system_user_id()`, this is consistent with how self-registration already bypasses the authenticated-context wrapper for the user insert itself.

One clarifying question worth resolving before you implement: do you want **strict `<`** (a level-100 user can never create another level-100 user) or **`<=` excluding self-level-equals-target** semantics for some future "manager" tier? Strict `<` is the safer default and is what I'd recommend unless you have a specific case for peer-creation at a non-admin level.

---

**what do you think about the project structure? I'm trying to following the clean architectures but in my own way, I generaly don't like the whole of any clean architecture**

---

Healthy instinct — "clean architecture, my own way" is usually the right call. Dogmatic adherence to Uncle Bob's original prescription (four rigid rings, use-case interactors, strict dependency-rule ceremony) tends to produce more indirection than an ERP backend actually needs. What you have is closer to a pragmatic layered architecture with DDD-influenced boundaries, which is a legitimate and common pattern — Go's own ecosystem norms (`golang-standards/project-layout`, most production services) lean this way rather than toward textbook Clean Architecture. Here's my honest read.

## What's working well

**The dependency direction is correct**, which is the actual point of clean architecture — everything else is ceremony. Handlers depend on services, services depend on repositories, repositories depend on `Querier` (an interface, not a concrete driver). Your domain layer (VOs) has zero outward dependencies. That's the property that matters; you don't need named "rings" to have it.

**The DTO/domain-model split is doing real work**, not just adding a layer for its own sake. DTOs are wire format, models carry invariants via VOs, and the service layer is the explicit conversion seam (`NewUserFromRegisterDTO`). This is more disciplined than most "clean" Go codebases I'd expect to see, and it's paying for itself — you're not leaking `json` tags into your domain, and you're not leaking bcrypt/VO logic into your DTOs.

**The `Querier` abstraction is the right amount of abstraction.** It solves one concrete problem (repositories work in or out of a transaction) without inventing a generic "unit of work" framework you don't need yet.

**Sealed interfaces (`__internal()`)** are a distinctly Go-idiomatic solution to a problem Clean Architecture usually solves with heavier machinery (explicit factories, DI containers). Good adaptation.

## What's inconsistent or unfinished

**The reflection-DI vs. constructor-injection question is still open** — we discussed this moving toward constructors, but `main.go` in front of me still shows the old reflection-based `RepoRegister`/`ServiceRegister`/`Add` pattern. If you're planning the migration, that's fine, just flagging it's not yet reflected in what's here.

**`Controller`, `Service`, `Repository` marker interfaces add a layer of indirection without adding a layer of behavior.** Right now they're just `interface { __internal() }` — they exist purely to gate implementations. That's a reasonable use of the sealed-interface pattern, but it's worth being honest that it's doing less than it looks like it's doing. If constructor injection replaces the reflection registry, you may not need `Service`/`Repository` as marker interfaces at all anymore — they existed mainly to let `ServiceRegister`/`RepoRegister` type-check a heterogeneous slice. Once wiring is explicit, `UserService` and `UserRepository` (the real interfaces with real methods) may be sufficient on their own.

**Query file organization will not survive scale.** `db/queries/*.sql` + `queries.go` with one `//go:embed` var per file works fine at 2 queries. At 50+ queries across a growing number of domains, a single flat `queries` package with one giant `embed.go` becomes a file you're constantly merge-conflicting on. Worth moving to `db/queries/user/`, `db/queries/roles/` now, before the pain shows up, since the migration cost only grows.

**`internal/services/service.go` and `internal/repositories/repository.go` are trivial single-method marker files** — fine as-is, but a sign this is genuinely a "your own way" architecture and not textbook DDD tactical patterns (no aggregate roots, no explicit repository-per-aggregate discipline yet). That's consistent with what your notes call "DDD-adjacent," and I think that's an accurate, non-judgmental label — you're borrowing DDD's _vocabulary and boundary discipline_ without its _tactical ceremony_ (specifications, domain events, aggregate invariant enforcement). That's a legitimate stopping point, not a partial failure to be "real" DDD.

## The actual gap: no explicit domain layer directory

Your VOs (`internal/vos`) are the domain layer, but they're not named or positioned as such — they sit alongside `models`, `dtos`, `services` as siblings, rather than being visually/structurally called out as the innermost ring. This is a nitpick, not a bug, but if you want the structure to _communicate_ the dependency direction to a new reader at a glance, consider:

```
internal/
  domain/
    vos/
    models/
  dtos/
  repositories/
  services/
  controllers/
  middlewares/
```

Nesting `models` under `domain/` alongside `vos` makes the "this is the center, everything else depends on it" property visible in the folder tree, not just true in practice. Purely organizational — doesn't change any behavior — but it's cheap and makes onboarding easier.

## Where I'd push back if you asked me to "complete" the Clean Architecture

I wouldn't add:

- **Use-case interactor structs** (one struct per use case, e.g. `RegisterUserUseCase`) — your service methods already are the use cases; wrapping each in its own struct+interface is the kind of ceremony that made you say you don't like "the whole" of clean architecture, for good reason.
- **A generic repository interface** (`Repository[T]`) — you already learned the real lesson (transaction-agnostic via `Querier`); a generic base adds abstraction without removing duplication that matters.
- **Full hexagonal "ports and adapters" renaming** — your `Controller`/`Service`/`Repository` naming is clearer to a Go reader than `Port`/`Adapter` would be.

## Bottom line

You're already doing the part of clean architecture that has real payoff — dependency direction, layer boundaries enforced by interfaces, business rules living in the domain rather than scattered across handlers. The parts you're skipping (use-case objects, generic repositories, ports/adapters vocabulary) are the parts that mostly exist to satisfy the _pattern_, not to solve a problem you actually have. I'd keep going in this direction rather than "completing" it toward the textbook version.

---

**give me a very small and concise but clear explanation of DDD so I can use it as a guide to improve my project**

---

## Domain-Driven Design, in one page

**Core idea:** the business logic (the _domain_) is the most important part of the code. Everything else — HTTP, SQL, JSON — exists to serve it, not the other way around. Structure the code so the domain doesn't know any of that other stuff exists.

### The vocabulary you actually need

- **Entity** — has an identity that persists over time, even as its data changes. A `User` is a `User` before and after an email update, because it's still _that same row_. Identity = the ID, not the field values.
- **Value Object (VO)** — has no identity, defined entirely by its data, and is immutable. `Email("a@b.com")` and `Email("a@b.com")` are simply _equal_, not "the same instance." Your `PasswordHash` is already this. VOs are where you put invariants ("an email must look like X", "a password must be ≥8 bytes").
- **Aggregate** — a cluster of entities/VOs treated as one unit for consistency. One entity in the cluster is the **aggregate root** — the only object outside code is allowed to reference or modify directly. Example: `Role` + its assignment rules might be an aggregate; you don't reach into its internals from outside, you talk to the root.
- **Repository** — loads and saves _one aggregate root_ at a time, by identity. Not a generic CRUD wrapper — it speaks in domain terms ("get the User aggregate"), not table rows. You're already close to this.
- **Domain Service** — logic that doesn't naturally belong to one entity or VO (e.g., "can this user assign this role?" — spans two aggregates). Lives in the domain layer, not your application `Service` layer, even though the names collide confusingly across the industry.
- **Bounded Context** — a boundary around a specific model where words mean one specific thing. "Role" in your Auth context might mean something different from "Role" in a future Payroll context. Each context gets its own models — you don't force one universal `Role` struct to serve every part of the system.
- **Ubiquitous Language** — the same words in code, conversation, and documentation. If the business says "suspend a user," the method is `Suspend()`, not `SetStatus(2)`.

### The one rule that matters more than any of the above

**Invariants live inside the aggregate, enforced by the aggregate — never checked externally and then trusted.** If "a user can only assign a role with a lower level," that rule should be nearly impossible to bypass by construction, not just checked once in a service method that someone might forget to call next time.

### How to use this as a guide, concretely

1. Look at your current `models` — ask "which of these have real identity vs. which are just data bags?" Identity → entity. No identity, just rules → VO (you're already doing this well for `PasswordHash`).
2. Group entities that must change together, atomically, into aggregates. `User` + its roles might be one aggregate if "assign a role" must be atomic with "check the level rule." Give it one root, one repository.
3. Push cross-cutting rules (like the role-level check) into a domain service or the aggregate root itself — not scattered across your application `service.go` files as ad-hoc `if` checks.
4. Don't chase the vocabulary for its own sake. If a concept in your app has no real identity question, no real invariant, and no cross-aggregate rule — it's fine as a plain struct. Not everything needs to be an Entity or VO.

That's the whole thing. The rest of DDD literature (event storming, domain events, CQRS, sagas) is tooling for _applying_ these ideas at larger scale — skip it until you actually feel the pain those tools solve.

---

so, what do you think about the code? any recommendation?

---

A few observations on `internal/codecs`, plus one related issue in the middleware that uses it:

**1. Context keys are untyped strings**
`codecKey = "codec"` (and `loggerKey`, `requestIDKey` elsewhere) are plain strings. This risks collisions if any dependency also uses `"codec"` as a context key, and `go vet` won't catch it. Standard fix — unexported custom type:

```go
type contextKey string
const codecKey contextKey = "codec"
```

**2. `GetCodec` panics if the codec isn't set**

```go
func GetCodec(ctx context.Context) codecs.Codec {
	return ctx.Value(codecKey).(codecs.Codec)
}
```

Fine as long as `Codec` middleware always runs first, but it's a silent contract — nothing enforces ordering. Consider a comma-ok return, or at least a comment noting the dependency.

**3. `Content-Type` gates both decode _and_ encode**
The middleware rejects any request without `Content-Type`, then uses that same codec to encode the response. Two consequences:

- A `GET` with no body will get a `400` even if it doesn't need one — fine today since you only have `POST` routes, but it'll bite you as soon as you add reads.
- Request format and response format are conflated. A client sending `application/json` has no way to request an `application/xml` response (or vice versa) — that's normally the job of `Accept`, separate from `Content-Type`.

If XML support is genuinely needed for both directions, splitting decode-codec (from `Content-Type`) and encode-codec (from `Accept`, defaulting to JSON) would be more correct. If XML isn't actually a real requirement, dropping it and the registry indirection entirely would simplify the codebase.

**4. Registry isn't extensible**
`registry` is a private package-level map — any new content type means editing `codec.go` directly. Not urgent at 2 codecs, but if this grows, a `Register(contentType string, codec Codec)` function would let codecs live in their own files/packages.

**5. The `TODO` for missing `Content-Type`**

```go
if raw == "" {
    // TODO: add missing content type error
    w.WriteHeader(http.StatusBadRequest)
    return
}
```

Worth resolving before this ships — right now the client gets an empty body with a 400 and no indication why.

None of this is urgent — the interface design (sealed via `__internal()`, clean `Decode`/`Encode` contract) is solid and consistent with the rest of your layering. The context-key typing and the Content-Type/Accept conflation are the two I'd actually fix soon, since the second one will force a breaking change later once you add `GET` endpoints.

Want me to look at the middleware chain or the `services`/`repositories` layers next, or focus on the codecs fix first?

---

**give me an explanation about access token and refresh token**

---


Here's a JWT implementation fitting your layered architecture.

## Claims and token generation

```go
// internal/vos/token_vos.go
package vos

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("expired token")
)

type AccessClaims struct {
	jwt.RegisteredClaims
	UserID uuid.UUID `json:"uid"`
}

const (
	accessTokenTTL  = 15 * time.Minute
	refreshTokenTTL = 7 * 24 * time.Hour
)

func NewAccessToken(userID uuid.UUID, secret []byte) (signed string, err error) {
	claims := AccessClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(accessTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		UserID: userID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

func ParseAccessToken(raw string, secret []byte) (claims AccessClaims, err error) {
	token, err := jwt.ParseWithClaims(raw, &claims, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return secret, nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			err = ErrExpiredToken
			return
		}
		err = ErrInvalidToken
		return
	}
	if !token.Valid {
		err = ErrInvalidToken
	}
	return
}
```

Note: the refresh token is **not** a JWT here — it's a random opaque string, stored server-side so it can be revoked. This is the recommendation from before.

```go
// internal/vos/refresh_token_vos.go
package vos

import (
	"crypto/rand"
	"encoding/base64"
)

type RefreshToken string

func NewRefreshToken() (RefreshToken, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return RefreshToken(base64.RawURLEncoding.EncodeToString(b)), nil
}
```

## Model + repository

```go
// internal/models/refresh_token_model.go
package models

import (
	"time"

	"github.com/google/uuid"
)

type RefreshToken struct {
	ID        uuid.UUID  `db:"id"`
	UserID    uuid.UUID  `db:"user_id"`
	TokenHash string     `db:"token_hash"` // store a hash, not the raw token
	ExpiresAt time.Time  `db:"expires_at"`
	RevokedAt *time.Time `db:"revoked_at"`
}
```

Store a SHA-256 hash of the refresh token, not the raw value — same principle as password hashing, so a DB leak doesn't hand out valid sessions directly.

## Service layer

```go
// internal/services/auth_service.go
package services

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"

	"erp/internal/dtos"
	"erp/internal/middlewares"
	"erp/internal/models"
	"erp/internal/repositories"
	"erp/internal/vos"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrRefreshTokenInvalid = errors.New("invalid or expired refresh token")
)

type AuthService interface {
	Service
	Login(ctx context.Context, in dtos.UserLogin) (out dtos.TokenPair, err error)
	Refresh(ctx context.Context, rawRefresh string) (out dtos.TokenPair, err error)
	Logout(ctx context.Context, rawRefresh string) (err error)
}

type authServiceImpl struct {
	db            *pgxpool.Pool
	users         repositories.UserRepository
	refreshTokens repositories.RefreshTokenRepository
	jwtSecret     []byte
}

func (service authServiceImpl) __internal() {}

func hashToken(raw string) string {
	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:])
}

func (service authServiceImpl) issuePair(ctx context.Context, tx pgx.Tx, userID uuid.UUID) (out dtos.TokenPair, err error) {
	access, err := vos.NewAccessToken(userID, service.jwtSecret)
	if err != nil {
		return
	}

	rawRefresh, err := vos.NewRefreshToken()
	if err != nil {
		return
	}

	_, err = service.refreshTokens.Create(ctx, tx, models.RefreshToken{
		UserID:    userID,
		TokenHash: hashToken(string(rawRefresh)),
		ExpiresAt: time.Now().Add(vos.RefreshTokenTTL),
	})
	if err != nil {
		return
	}

	out = dtos.TokenPair{AccessToken: access, RefreshToken: string(rawRefresh)}
	return
}

func (service authServiceImpl) Login(ctx context.Context, in dtos.UserLogin) (out dtos.TokenPair, err error) {
	logger := middlewares.GetLogger(ctx)
	model := models.NewUserFromLoginDTO(in)

	err = withTx(ctx, service.db, func(tx pgx.Tx) (err error) {
		user, err := service.users.GetByName(ctx, tx, model)
		if err != nil {
			return err
		}
		if !user.PasswordHash.Matches(in.Password) {
			return ErrIncorrectUserPassword
		}

		out, err = service.issuePair(ctx, tx, user.ID)
		return err
	})
	if err != nil {
		logger.ErrorContext(ctx, "login failed", "error", err)
	}
	return
}

func (service authServiceImpl) Refresh(ctx context.Context, rawRefresh string) (out dtos.TokenPair, err error) {
	hash := hashToken(rawRefresh)

	err = withTx(ctx, service.db, func(tx pgx.Tx) (err error) {
		stored, err := service.refreshTokens.GetValidByHash(ctx, tx, hash)
		if errors.Is(err, repositories.ErrNotFound) {
			return ErrRefreshTokenInvalid
		} else if err != nil {
			return err
		}

		// rotation: revoke old, issue new
		if err = service.refreshTokens.Revoke(ctx, tx, stored.ID); err != nil {
			return err
		}

		out, err = service.issuePair(ctx, tx, stored.UserID)
		return err
	})
	return
}

func (service authServiceImpl) Logout(ctx context.Context, rawRefresh string) (err error) {
	hash := hashToken(rawRefresh)
	return withTx(ctx, service.db, func(tx pgx.Tx) error {
		return service.refreshTokens.RevokeByHash(ctx, tx, hash)
	})
}
```

## Middleware for verifying access tokens

```go
// internal/middlewares/auth_middleware.go
package middlewares

import (
	"context"
	"net/http"

	"erp/internal/vos"

	"github.com/google/uuid"
)

const userIDKey = "user-id"

func Auth(secret []byte) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			raw := r.Header.Get("Authorization")
			raw = strings.TrimPrefix(raw, "Bearer ")
			if raw == "" {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			claims, err := vos.ParseAccessToken(raw, secret)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), userIDKey, claims.UserID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetUserID(ctx context.Context) uuid.UUID {
	return ctx.Value(userIDKey).(uuid.UUID)
}
```

## Key design points

- **Access token is stateless**: `ParseAccessToken` never hits the DB, so verification stays fast on every request.
- **Refresh token is opaque and stored hashed**: gives you a revocation point without weakening the access-token performance benefit.
- **Rotation on refresh**: each `Refresh` call invalidates the old refresh token and issues a new one. If a revoked token is reused, that's a strong signal of theft — worth revoking the entire session family (all tokens tied to that user/device) rather than just logging the anomaly.
- **`GetUserID`** in middleware mirrors your existing `GetLogger`/`GetCodec` context pattern, keeping the auth middleware consistent with the rest of your stdlib-based design.

Want the `RefreshTokenRepository` implementation and the corresponding migration file next?


---

**dame una guia de como usar el modulo jwt**

---

## Guía práctica: golang-jwt/jwt/v5

### Instalación

```bash
go get github.com/golang-jwt/jwt/v5
```

### Conceptos básicos

El módulo trabaja con tres piezas:
- **Claims**: los datos que va dentro del token (quién es el usuario, cuándo expira, etc.)
- **Signing method**: el algoritmo de firma (HMAC, RSA, ECDSA...)
- **Token**: la estructura que combina claims + firma

### 1. Definir tus claims

Puedes usar `jwt.RegisteredClaims` (los campos estándar del spec: `exp`, `iat`, `sub`, etc.) y extenderlo con tus propios campos:

```go
type AccessClaims struct {
	jwt.RegisteredClaims
	UserID uuid.UUID `json:"uid"`
	Role   string    `json:"role,omitempty"`
}
```

Los campos estándar más usados de `RegisteredClaims`:

| Campo | JSON | Propósito |
|---|---|---|
| `ExpiresAt` | `exp` | Cuándo expira |
| `IssuedAt` | `iat` | Cuándo se emitió |
| `NotBefore` | `nbf` | No válido antes de esta fecha |
| `Subject` | `sub` | Identificador del sujeto (opcional, si no usas campo custom) |
| `ID` | `jti` | ID único del token (útil para blocklists) |

### 2. Firmar (generar) un token

```go
claims := AccessClaims{
	RegisteredClaims: jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	},
	UserID: userID,
}

token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
signed, err := token.SignedString([]byte(secretKey))
```

`SignedString` acepta distintos tipos según el algoritmo:
- HMAC (`HS256/384/512`) → `[]byte`
- RSA (`RS256`...) → `*rsa.PrivateKey`
- ECDSA (`ES256`...) → `*ecdsa.PrivateKey`

Para tu caso (monolito, un solo servicio verificando), **HMAC es suficiente y más simple** — no necesitas el par de llaves pública/privada que justifica RSA/ECDSA (eso tiene sentido cuando otro servicio necesita *verificar* sin poder *firmar*).

### 3. Parsear y validar un token

```go
func ParseAccessToken(raw string, secret []byte) (claims AccessClaims, err error) {
	token, err := jwt.ParseWithClaims(raw, &claims, func(t *jwt.Token) (any, error) {
		// valida que el algoritmo sea el esperado (evita ataques de "alg confusion")
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return secret, nil
	})
	if err != nil {
		return
	}
	if !token.Valid {
		err = errors.New("invalid token")
	}
	return
}
```

Puntos clave:
- La función callback (`Keyfunc`) recibe el token *antes* de validarlo, para que decidas qué llave usar según el `alg` — por eso **siempre** se verifica el signing method dentro del callback, nunca se asume.
- `jwt.ParseWithClaims` valida automáticamente `exp` y `nbf` si están presentes.

### 4. Manejar errores específicos

```go
token, err := jwt.ParseWithClaims(raw, &claims, keyFunc)
if err != nil {
	switch {
	case errors.Is(err, jwt.ErrTokenExpired):
		// el token expiró — el cliente debe usar el refresh token
	case errors.Is(err, jwt.ErrTokenMalformed):
		// el token está corrupto o mal formado
	case errors.Is(err, jwt.ErrTokenSignatureInvalid):
		// la firma no coincide — posible manipulación
	default:
		// otro error de validación
	}
}
```

Distinguir `ErrTokenExpired` del resto es importante en tu flujo: un token expirado es un caso normal (activa el refresh), mientras que una firma inválida es sospechoso y debería loguearse como tal.

### 5. Opciones de validación adicionales

```go
token, err := jwt.ParseWithClaims(raw, &claims, keyFunc,
	jwt.WithValidMethods([]string{"HS256"}),        // refuerza el algoritmo permitido
	jwt.WithLeeway(30*time.Second),                  // tolerancia de reloj
	jwt.WithIssuer("erp-api"),                       // valida el campo iss
)
```

`WithValidMethods` es una capa extra de seguridad (además del check dentro del `Keyfunc`) — recomendable si vas a exponer el sistema a más de un tipo de cliente en el futuro.

### 6. Aplicado a tu arquitectura

En tu `vos` package esto encapsula bien la invariante "un access token válido siempre tiene un `UserID` y no está expirado":

```go
// internal/vos/token_vos.go — value object, sin lógica de negocio externa
func NewAccessToken(userID uuid.UUID, secret []byte) (string, error) { ... }
func ParseAccessToken(raw string, secret []byte) (AccessClaims, error) { ... }
```

Y la capa de `middlewares` es quien lo consume, siguiendo el mismo patrón que ya usas con `GetLogger`/`GetCodec`:

```go
claims, err := vos.ParseAccessToken(raw, secret)
if err != nil { /* 401 */ }
ctx := context.WithValue(r.Context(), userIDKey, claims.UserID)
```

### Resumen de buenas prácticas

- Nunca confíes en el `alg` del header sin validarlo explícitamente en el `Keyfunc`
- Guarda el `secret` en config (vos ya usan `caarlos0/env`), nunca hardcodeado
- Usa `RegisteredClaims` como base — te da `exp`/`iat` gratis con la validación correcta
- Para el refresh token, no uses JWT (como vimos antes) — usa un token opaco almacenado en DB, para poder revocarlo

¿Quieres que arme el ejemplo del middleware completo con extracción del header `Authorization` y manejo de los distintos casos de error?

---

**que es keyfunc**

---

## Keyfunc

`Keyfunc` es el tipo de la función callback que le pasas a `jwt.ParseWithClaims` (o `jwt.Parse`) para que le diga a la librería **qué llave usar para verificar la firma** del token.

```go
type Keyfunc func(*Token) (interface{}, error)
```

### Por qué existe

La librería no puede saber de antemano qué llave usar — depende del token. Entonces el flujo es:

1. La librería decodifica el **header** del JWT (sin validar nada todavía)
2. Te pasa el `*jwt.Token` parcialmente decodificado a tu `Keyfunc`
3. Tu función mira ese token (típicamente el campo `alg`, o a veces `kid` — key ID) y decide qué llave devolver
4. La librería usa esa llave para verificar la firma

```go
keyFunc := func(t *jwt.Token) (interface{}, error) {
    // t.Method te dice qué algoritmo dice el header que usó
    if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
        return nil, errors.New("unexpected signing method")
    }
    return secret, nil // esta es la llave que se usará para verificar
}

token, err := jwt.ParseWithClaims(raw, &claims, keyFunc)
```

### Por qué siempre hay que validar el algoritmo dentro

Esto es el punto más importante del `Keyfunc`, no un detalle opcional. Existe un ataque clásico llamado **"algorithm confusion"**: si tu API pública fuera RSA (`RS256`) y firmaras el token con la llave privada, un atacante podría:

1. Tomar el token válido
2. Cambiar el header a `alg: HS256`
3. Firmarlo con la llave *pública* (que es conocida) como si fuera un secreto HMAC

Si tu `Keyfunc` no verifica que el algoritmo sea el esperado antes de devolver la llave, la librería aceptaría esa firma falsa como válida — porque técnicamente "coincide" con lo que le diste.

Por eso el patrón correcto es siempre:

```go
func(t *jwt.Token) (interface{}, error) {
    if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
        return nil, errors.New("unexpected signing method")
    }
    return secret, nil
}
```

En vez de simplemente:

```go
func(t *jwt.Token) (interface{}, error) {
    return secret, nil // ❌ nunca verifica el alg — vulnerable
}
```

### Cuándo el Keyfunc necesita más lógica

Si en algún momento rotas llaves de firma (por ejemplo, cambias el secreto cada cierto tiempo por seguridad), el `Keyfunc` es también el lugar donde resolverías eso, normalmente usando el claim `kid` (key ID) del header para buscar la llave correcta en un mapa o almacén de llaves:

```go
func(t *jwt.Token) (interface{}, error) {
    kid, ok := t.Header["kid"].(string)
    if !ok {
        return nil, errors.New("missing kid")
    }
    key, ok := keyStore[kid]
    if !ok {
        return nil, errors.New("unknown key")
    }
    return key, nil
}
```

Para tu caso actual (un solo secreto HMAC, sin rotación), la versión simple con el chequeo de `alg` es suficiente.
