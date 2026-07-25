# ADR-0003: HTTP Router — stdlib `net/http.ServeMux`

## Status

Accepted

## Context

The project needs an HTTP router for request dispatch. Go 1.22 introduced
enhanced pattern routing (method matching, path wildcards) to the standard
library's `ServeMux`, narrowing the gap with third-party routers.

## Decision

Use the standard library's `net/http.ServeMux`, with a thin custom
abstraction layered on top for middleware chaining and route grouping, to be
implemented as the project's routing needs require it.

## Alternatives considered

`chi` — evaluated informally as the benchmark for a minimalist third-party
router. Notably, even chi's own maintainers do not migrate it to build on top
of the new stdlib `ServeMux` internally, which was read as a signal that the
stdlib router is now capable enough for most use cases without a wrapping
dependency.

## Consequences

- (+) Zero third-party dependency for routing; smaller supply-chain surface.
- (+) One less abstraction to learn/maintain; behavior fully documented in the
  stdlib.
- (−) No native middleware chaining or route grouping — requires a thin
  custom abstraction on top.
- (−) No native route introspection/listing for debugging or auto-generated
  documentation.
- (−) Path parameter access via `r.PathValue()` is stringly-typed, with no
  struct-binding support out of the box.
