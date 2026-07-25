# ADR-0005: Audit Columns Are Database-Only, Excluded From Service-Layer Models

## Status

Accepted

## Context

Every persisted table carries audit metadata — `created_at`, `updated_at`,
`deleted_at`, `created_by_id`, `updated_by_id`, `deleted_by_id` — for
internal traceability. The question is whether this metadata should be part
of the domain/service-layer representation of an entity, or remain confined
to the persistence layer.

## Decision

Audit columns exist only in the database schema and in the repository
layer's persistence-mapping structs. They are not part of the models that
services operate on or that get exposed through the service/handler layers.

## Rationale

Audit metadata answers "who touched this row and when," which is an
infrastructure/traceability concern, not a domain concern. Including it in
service-layer models would leak persistence details upward and invite
domain logic to depend on fields that have no business meaning (e.g., a
service method branching on `updated_at`). Keeping it repository-scoped
preserves the layering discipline: infrastructure detail stays behind the
repository boundary.

## Alternatives considered

- **Embed audit fields in domain models** (e.g., an embedded `AuditFields`
  struct on every entity): rejected — couples domain models to a
  persistence concern and makes it available (and temptingly usable) to
  business logic that should not depend on it.

## Consequences

- (+) Domain and service-layer models stay focused on business-relevant
  fields only.
- (+) Audit columns can evolve (e.g., adding a new tracked action) without
  touching service-layer code.
- (−) If a future feature genuinely needs audit data as business information
  (e.g., displaying "last edited by" in a UI), it must be deliberately
  exposed through a dedicated read path rather than reused from the
  persistence struct.
