# ADR-0007: Audit User Attribution — `SET LOCAL app.current_user_id`

## Status

Accepted

## Context

Audit columns (`created_by_id`, `updated_by_id`, `deleted_by_id`) need to be
populated automatically at the database layer (see ADR-0005 and ADR-0004),
but the database has no inherent notion of "which application user issued
this statement" — that identity lives in the application's request context,
not in the SQL connection.

## Decision

At the start of each transaction, the application issues
`SET LOCAL app.current_user_id = '<user-id>'`. This happens in exactly one
place: the shared transaction-opening helper that every repository already
goes through to get automatic commit/rollback handling — no individual
repository method sets it directly.

Triggers responsible for populating audit columns (including the
`INSTEAD OF DELETE` trigger from ADR-0004) read the actor through a
dedicated database function instead of calling `current_setting` directly.
That function returns the value of `app.current_user_id` for the current
transaction when it's set; if it isn't set (e.g., a background job or a
direct `psql` session), it falls back to a reserved system user ID instead
of leaving the audit column null.

## Rationale

`SET LOCAL` scopes the setting to the current transaction only — it's
automatically reset on commit or rollback, so there's no risk of a stale
user ID leaking into a pooled connection's next transaction. Setting it
inside the shared transaction helper (rather than in each repository
method) means there is exactly one place where it can be forgotten, not one
per write path — and even then, the fallback-to-system-user behavior keeps
the database usable and audit columns non-null. This keeps attribution
enforcement at the database layer (consistent with ADR-0004's approach of
pushing invariants down to the schema) instead of requiring every service
method to manually pass and set an actor ID on every write.

## Alternatives considered

- **Pass `created_by_id`/`updated_by_id` explicitly as columns on every
  `INSERT`/`UPDATE` from the service layer**: rejected — reintroduces the
  discipline problem ADR-0004 and ADR-0005 are meant to avoid; a forgotten
  parameter silently attributes a change to the wrong user (or none).

## Consequences

- (+) Attribution is enforced at the database layer; application code
  cannot forget to stamp an actor.
- (+) Automatically scoped and cleared per transaction — no manual cleanup.
- (+) Centralized in the shared transaction helper, so there is a single
  place to get right instead of one per repository method.
- (+) The fallback-to-system-user behavior means an unset actor (e.g., a
  background job or a raw `psql` session) degrades gracefully to a
  reserved system user rather than leaving audit columns null.
- (−) Any code path that opens a transaction without going through the
  shared helper bypasses attribution entirely and silently falls back to
  the system user — this makes bypassing the helper a real risk to guard
  against, not just an inconvenience.
