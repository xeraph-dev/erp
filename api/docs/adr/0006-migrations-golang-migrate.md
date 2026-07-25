# ADR-0006: Schema Migration Tool — golang-migrate

## Status

Accepted

## Context

The project needs a tool to manage versioned, repeatable PostgreSQL schema
migrations as part of the development and deployment workflow.

## Decision

Use `golang-migrate`.

## Rationale

Prior working proficiency; it is straightforward and comfortable to use in
practice — plain up/down SQL files, a simple CLI, and a Go library that
integrates cleanly with the existing `pgx`-based stack without imposing an
ORM or a custom migration DSL.

## Alternatives considered

None formally evaluated — the choice was made directly based on existing
familiarity, not a comparative evaluation.

## Consequences

- (+) Low ramp-up time; plain SQL migrations are easy to review and reason
  about.
- (+) No additional abstraction (DSL, ORM) introduced on top of raw SQL.
- (−) No comparative data on how alternatives (e.g., `goose`, `atlas`) would
  have performed for this use case; this decision cannot be revisited
  against concrete trade-off data later, only against future needs.
