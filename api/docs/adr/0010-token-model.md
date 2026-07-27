# ADR-0010: Token Model — Short-Lived JWT Access Tokens, Opaque Long-Lived Refresh Tokens

## Status

Accepted

## Context

Every authenticated request needs to verify the caller's identity, and the
verification mechanism has two competing concerns: it should be cheap
(ideally no database round trip on the common path), and it should be
revocable (a compromised or logged-out session must stop working promptly).
A single token type struggles to satisfy both — a self-contained token that
needs no database lookup to verify is, by the same property, hard to revoke
before it expires; a database-backed token that can be revoked immediately
requires a lookup on every use.

## Decision

Use two distinct token types, each optimized for one of the two concerns:

- **Access token**: a JWT (`HS256`), carrying the user ID as `sub` and
  standard `iat`/`nbf`/`exp` claims, with a short TTL (~15 minutes). It is
  verified purely by signature and claim checks — no database access is
  required to authenticate a request.
- **Refresh token**: an opaque random string (not a JWT), stored hashed in
  PostgreSQL under the append-only family model (see ADR-0008), with a long
  TTL. Since it carries no information of its own, revoking it is a matter
  of updating its database row.

Logging in or refreshing issues both tokens as a pair; only the refresh
token is persisted server-side.

## Rationale

The access token's short TTL bounds the exposure window of a stolen token
without needing a revocation mechanism for it: even if invalidating it
immediately isn't possible, waiting fifteen minutes for it to expire on its
own is an acceptable risk trade-off in exchange for skipping a database
round trip on every request. This keeps steady-state authenticated request
latency low and avoids adding database load proportional to request volume,
which matters more than instant revocation for a token whose blast radius
is already time-boxed.

The refresh token needs the opposite property — it must be revocable
immediately, since logout and reuse detection (ADR-0008) both depend on a
revoked token being rejected on its very next use. Making it opaque rather
than a self-contained JWT means it carries no claims to trust or distrust;
its only source of truth is the database row, which is exactly what allows
it to be revoked with a simple, strongly-consistent update.

Access tokens are not individually revocable today; only their natural
expiry ends them. Should individual-token revocation or per-token
traceability become necessary later, a `jti` claim can be added at that
point — it isn't included now to avoid carrying a claim with no current
consumer.

## Alternatives considered

- **A single long-lived JWT for everything**: rejected — no revocation path
  short of a server-side blocklist, which reintroduces a database check on
  every request and defeats the reason for choosing a JWT in the first
  place.
- **JWTs for both access and refresh tokens**: rejected — reuse detection
  and logout (ADR-0008) require checking and mutating revocation state on
  every refresh, so the refresh token would need a blocklist regardless of
  being a JWT; an opaque, database-backed token gets the same guarantee
  more directly, without carrying unused claims.
- **A single long-lived opaque session token (no JWT at all)**: rejected —
  every authenticated request would require a database round trip to
  verify identity, adding latency and load that scales with request volume
  rather than login volume.

## Consequences

- (+) Authenticating a request within the access token's TTL requires no
  database access.
- (+) A stolen access token's usable window is bounded to ~15 minutes even
  without an active revocation check.
- (+) Refresh tokens are immediately revocable through the family model in
  ADR-0008, satisfying logout and reuse-detection requirements.
- (−) Two separate verification paths exist (signature/claims for access
  tokens, database lookup for refresh tokens), each needing its own
  correctness care.
- (−) An individual access token cannot be revoked before it expires; a
  compromised access token remains valid for up to its full TTL regardless
  of server-side action.
- (−) Clock skew between issuer and future distributed verifiers (should
  the system ever run more than one instance) requires tolerance in JWT
  validation (`jwt.WithLeeway`), which slightly widens the token's
  effective validity window.
