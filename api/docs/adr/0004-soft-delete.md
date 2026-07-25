# ADR-0004: Soft Delete — Exclusion via Postgres Views

## Status

Accepted

## Context

The system needs a soft-delete mechanism (marking records as deleted without
physically removing them, for audit and recovery purposes) that is
consistently enforced without relying on every repository to remember to
filter deleted rows manually.

## Decision

Use a `deleted_at` column per soft-deletable table, combined with a Postgres
view per table that excludes rows where `deleted_at IS NOT NULL`. Application
code reads and writes through the view, not the underlying table, for all
standard operations.

Because each view selects from a single base table with no joins or
aggregation, Postgres treats it as an automatically updatable view: `INSERT`
and `UPDATE` against the view translate directly to the base table with no
extra plumbing. A `DELETE` against the view is intercepted by an
`INSTEAD OF DELETE` trigger that rewrites it into
`UPDATE ... SET deleted_at = now()` on the base table, so callers can issue
a normal `DELETE` without needing to know it's a soft delete under the hood.

## Rationale

Enforcing the filter at the database layer (via a view) removes the
possibility of a repository forgetting a `WHERE deleted_at IS NULL` clause
and leaking soft-deleted rows. It keeps the exclusion rule in one place
(the view definition) instead of duplicated across every query.

## Alternatives considered

- **`WHERE deleted_at IS NULL` in every repository query**: rejected —
  relies on developer discipline in every single query; a missed clause is a
  silent data leak.
- **Row-Level Security (RLS) policies**: not chosen — adds session-level
  configuration overhead (setting policy context per connection) that isn't
  justified for a single-tenant application; views achieve the same
  exclusion guarantee with less operational complexity.

## Consequences

- (+) Soft-delete filtering is enforced once, at the schema level — no
  per-query discipline required.
- (+) The view provides a uniform CRUD surface (`SELECT`/`INSERT`/`UPDATE`/
  `DELETE`) that behaves like a normal table to application code, while
  transparently turning deletes into soft deletes.
- (−) Every soft-deletable table requires a paired view and an
  `INSTEAD OF DELETE` trigger, both kept in sync via migrations.
- (−) Direct access to soft-deleted rows (e.g., for administrative recovery)
  requires explicitly querying the base table, bypassing the view.
