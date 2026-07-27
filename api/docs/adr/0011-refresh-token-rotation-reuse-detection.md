# ADR-0011: Refresh Token Rotation & Reuse Detection — Strict, No Grace Window

## Status

Accepted

## Context

ADR-0008 decided *where* refresh tokens live (PostgreSQL, hashed, in
append-only families) and named rotation and reuse detection as the
capabilities that storage model needs to support. This ADR covers the
*policy* built on top of that storage: what happens, step by step, when a
refresh token is used, and how strictly reuse is treated.

Rotating a refresh token on every use is what makes reuse detectable in the
first place — if a token were reusable, there would be nothing to notice
when an attacker replays a stolen one. But rotation introduces its own
failure mode: two requests racing to refresh the same token concurrently
(a real scenario with a client that retries, or double-fires a request)
must not both succeed, or the token has effectively not been single-use.

## Decision

On every refresh:

1. The presented token is looked up by its hash.
2. If it is already marked `revoked_at`, this is treated as **reuse**, and
   the *entire family* is revoked immediately — every token that ever
   descended from the same login. The request is rejected.
3. If it is expired, the request is rejected without revoking the family.
4. Otherwise, the token is revoked and a new token is appended to the same
   family in a single transaction, using `SELECT FOR UPDATE` to lock the
   row being rotated. This closes the race window: a concurrent refresh
   attempt using the same token blocks until the first transaction commits,
   then sees the token as already revoked and is treated as reuse per step
   2.

Logging out revokes the entire family directly, without requiring reuse to
be detected first. `RevokeFamily` is the single shared primitive behind
both paths — reuse detection and logout are, mechanically, the same
operation triggered by different conditions.

There is no grace window anywhere in this flow: a token is either valid for
exactly one rotation or it is treated as compromised. No brief overlap
period is allowed where both the old and new token work, even though such
an overlap is a common pattern to smooth over client-side race conditions
(e.g., a mobile app with two tabs).

## Rationale

Treating any reuse of a rotated-out token as theft — rather than as a
possible benign race — is a deliberate strict stance: a legitimate client
never has a reason to present a token it has already exchanged for a new
one, since it received and should be using the replacement. The only
realistic explanation for a rotated-out token being replayed is that an
attacker captured it independently of the legitimate client's own rotation.
Reacting to that signal by revoking the whole family, not just the reused
token, is what makes the mechanism worth having: an attacker who is one
step behind the legitimate client's rotation gets cut off the moment they
try to use their copy, instead of continuing to sit on a still-valid stolen
token.

The row lock during rotation exists because "reject reuse" and "handle
concurrent legitimate requests" would otherwise be indistinguishable: two
near-simultaneous refresh calls using the same valid token are, from the
database's point of view, identical to one legitimate call and one replay,
unless something serializes them. Locking forces one to complete first, so
the second one deterministically sees an already-revoked token and is
handled by the same reuse-detection path — no separate "concurrent request"
case needs to be reasoned about.

Rejecting a grace window was a conscious trade-off: it removes an entire
class of edge cases (how long is the grace period, does it reopen the
theft window, does it complicate the lock) at the cost of being less
forgiving of legitimate client-side races that aren't handled by the
locking — for example, multiple browser tabs or client instances sharing
the same refresh token, where one tab rotates it and another, unaware of
the rotation, tries to use the now-stale copy moments later. That failure
mode is treated as a client-side concern to design around (e.g., cross-tab
coordination), not a case to accommodate in the server-side policy.

## Alternatives considered

- **Grace window on rotation** (old token still valid for a few seconds
  after rotation): rejected — reintroduces a window where a stolen token
  is usable even after the legitimate client has rotated past it, directly
  undermining the reason for rotating in the first place.
- **Revoke only the reused token, not the whole family**: rejected — an
  attacker holding a stolen token can simply request a fresh rotation from
  their copy and continue undetected; the legitimate user's session is
  compromised without either party finding out. Revoking the whole family
  forces both the attacker and the legitimate user to re-authenticate,
  which is the correct outcome once theft is suspected.
- **Optimistic concurrency (compare-and-swap on the token row) instead of
  row locking**: rejected — would require the losing concurrent request to
  retry or fail, and distinguishing "lost the race legitimately" from "this
  is reuse" adds complexity without a clear benefit over blocking briefly
  via `SELECT FOR UPDATE`, which is simple and already sufficient at this
  project's request volume.

## Consequences

- (+) A stolen refresh token is only usable until the legitimate client's
  next rotation, and using it after that point immediately invalidates the
  entire session lineage for both parties.
- (+) Concurrent refresh requests using the same token are handled by the
  same code path as reuse detection — no separate race-handling logic.
- (+) `RevokeFamily` being shared between logout and reuse detection keeps
  the revocation mechanism in one place.
- (−) Legitimate but poorly-behaved clients (e.g., firing duplicate,
  non-retried refresh requests) will trip reuse detection and be logged
  out, even though no theft occurred. This pushes the burden of not
  double-refreshing onto client implementations.
- (−) Row locking during rotation serializes concurrent refresh attempts
  for the same token, adding minor latency to the losing request (though
  it is expected to fail anyway).
