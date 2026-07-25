# ADR-0001: Backend Language — Go

## Status

Accepted

## Context

The project needs a backend language for a modular monolith ERP system,
prioritizing long-term maintainability, performance under concurrent load,
and simple deployment (single binary).

## Decision

Use Go.

## Rationale

Prior working proficiency with the language removed ramp-up cost. Go's
relevant qualities for this use case:

- Built-in concurrency primitives (goroutines/channels).
- Static typing.
- Fast compilation.
- A single statically-linked binary for deployment.
- A strong standard library (`net/http`, `database/sql`) that reduces
  dependency surface.

## Alternatives considered

None formally evaluated — the choice was made directly based on existing
familiarity, not a comparative evaluation.

## Consequences

- (+) Low ramp-up time; deployment simplicity (single binary embedding the
  frontend).
- (+) Strong fit for a modular monolith with clear layer boundaries.
- (−) No comparative data exists on how alternatives (e.g., Node.js/TypeScript,
  Java) would have performed for this use case; this decision cannot be
  revisited against concrete trade-off data later, only against future needs.
