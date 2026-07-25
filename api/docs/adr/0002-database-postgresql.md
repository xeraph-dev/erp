# ADR-0002: Primary Datastore — PostgreSQL

## Status

Accepted

## Context

An ERP system's core domain (users, roles, transactions, relational business
data) requires strong consistency guarantees and relational integrity.

## Decision

Use PostgreSQL as the primary datastore.

## Rationale

Prior working proficiency; ACID compliance and mature relational modeling are
a strong fit for ERP domain data. Native support for constraints, triggers,
and views (already used for soft-delete patterns and audit columns in this
project) covers domain rules directly at the data layer.

## Alternatives considered

None formally evaluated — chosen directly based on existing familiarity.

## Consequences

- (+) Relational integrity enforced at the DB layer, reducing bugs from
  partial writes.
- (+) Mature tooling (`pgx`, `golang-migrate`) already integrated into the
  stack.
- (−) No comparative data on NoSQL/other relational alternatives for this
  specific workload.
