# ADR-0008: Refresh Token Storage — PostgreSQL, Append-Only Token Families

## Status

Accepted

## Context

Refresh tokens need a storage model that supports rotation on every use,
detection of token reuse (a strong signal of theft), and revocation of an
entire session lineage when reuse is detected or the user logs out. The
question is which datastore should hold this state, and how the rotation
model should be structured.

## Decision

Store refresh tokens in PostgreSQL, hashed at rest (the raw token is never
persisted — only its hash, so a database leak doesn't expose usable
tokens). Each token belongs to a family:

- Logging in creates a new family.
- Refreshing issues a new token in the same family and revokes the token
  that was just used.
- Logging out revokes the entire family.
- Attempting to reuse an already-revoked token revokes the entire family
  immediately, since reuse of a dead token is treated as evidence the
  token was stolen — the legitimate client would never present a token it
  already exchanged.

The family is an append-only chain: every rotation adds a new row rather
than overwriting the previous one, so the full lineage of a session remains
inspectable. Rotation (issuing the new token and revoking the old one) and
reuse detection both go through row locking within a transaction, to avoid
a race where two concurrent refresh requests using the same token both
succeed.

## Rationale

The rotation and reuse-detection model requires strong, immediate
consistency — a revoked token must be seen as revoked by the very next
request, since the whole point of reuse detection is catching theft as
soon as it happens. PostgreSQL gives this for free through ordinary
transactions and row locking. The family also has a natural relationship
to a user (a foreign key), which is a simple, integrity-checked relation in
a relational database.

Being append-only rather than mutating a single row per session keeps a
full audit trail of a session's rotation history, which is useful both for
debugging and for post-incident review if a theft is ever detected.

## Alternatives considered

- **Redis as the primary store**: rejected for now. Redis would need a
  Lua script (or `WATCH`/`MULTI`) to get the same atomicity that PostgreSQL
  transactions provide natively, and durability isn't the default behavior
  — it would need explicit persistence configuration, or a Redis restart
  would silently log out every active session. Redis remains a candidate
  for a future read-optimization layer in front of PostgreSQL, not as the
  system of record from the start.

## Consequences

- (+) Strong consistency for rotation and reuse detection, backed by
  PostgreSQL's transactional guarantees.
- (+) Full session lineage is queryable, useful for security review and
  future features like "view active sessions."
- (+) A leaked database does not expose usable refresh tokens, since only
  hashes are stored.
- (−) Every refresh request costs a database round trip (plus row locking),
  which is more expensive than an in-memory lookup — acceptable at this
  project's scale, but a candidate for revisiting if refresh volume grows
  significantly.
- (−) Expired or fully-revoked families need an explicit cleanup process
  (e.g., a periodic job), since PostgreSQL has no built-in TTL the way
  Redis does.
