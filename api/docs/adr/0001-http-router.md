# ADR-0001: HTTP Router

## Status
Accepted

## Context
The project needs an HTTP routing mechanism. Go 1.22+ introduced
method-aware pattern matching in net/http.ServeMux, closing most of
the gap with third-party routers.

## Decision
Use the native net/http.ServeMux, wrapped in a thin custom layer
(Router/Group) to support middleware chaining and route grouping.

## Alternatives considered
- chi: more expressive, but adds a dependency for marginal benefit
  given a custom Group/Chain layer is already implemented.
- gin/echo: full frameworks that impose their own Handler model,
  conflicting with the existing codec/DTO design.
- gorilla/mux: repository is no longer actively maintained.

## Consequences
- Zero external dependencies for routing.
- Ownership of maintaining Group/Chain falls on the project if
  requirements grow (e.g. more complex route params).
