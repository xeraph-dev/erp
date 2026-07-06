# Technology Stack & Infrastructure Overview

This document outlines the rationale behind the selected technologies and defines the high-level communication boundaries between system components.

---

## Core Technologies

| Component | Technology | Decision Rationale |
| :--- | :--- | :--- |
| **API / Backend** | Go | High concurrency (goroutines) for handling many simultaneous ERP requests. Strong standard library, strict typing, and compiled binaries ensure reliability and performance in production. |
| **Frontend / UI** | Next.js | Delivers a modern React-based user interface with server-side rendering (SSR) for complex dashboards, enabling good SEO (for public-facing parts) and optimal client-side performance. |
| **Database** | PostgreSQL | ACID compliance is non-negotiable for financial/ERP data. Strong support for Row-Level Security (RLS), JSONB flexibility, and robust transactional integrity. |
| **Cache & Queues** | Redis | Ultra-low latency for permission caching, distributed rate-limiting, and decoupled async job queues (e.g., audit logging). Acts as a stateless coordination layer across Go instances. |
| **Reverse Proxy** | Nginx | Terminates TLS/SSL, enforces global rate limits, serves static Next.js assets, and load-balances traffic to the Go API instances. |

---

## System Boundaries & Communication Flow

The system is structured into distinct layers to enforce separation of concerns.

**Request Flow:**

1. All external traffic (browser, mobile, API clients) arrives via **HTTPS** to **Nginx**.
2. Nginx performs:
   - TLS termination.
   - Global rate limiting (login, registration, password reset).
   - Static asset serving (for Next.js built files).
3. Nginx routes requests based on the path:
   - `/api/*` requests are proxied to the **Go API** layer.
   - UI routes (`/` and other frontend paths) are proxied to the **Next.js** server.
4. The **Go API** layer:
   - Executes business logic and authorization checks.
   - Reads/writes to **PostgreSQL** (source of truth).
   - Interacts with **Redis** for permission caching and background job queuing.
5. **Next.js**:
   - Renders server-side pages when needed.
   - Fetches data from the Go API for dynamic content.
   - Delegates client-side navigation to React.
6. **Redis** serves two roles:
   - **Cache:** Stores user permissions for low-latency authorization checks.
   - **Queue:** Buffers audit log entries for asynchronous batch writing.
7. **PostgreSQL**:
   - Enforces Row-Level Security (RLS) as a defense-in-depth measure.
   - Maintains referential integrity and transactional consistency.

---

## Architecture Principles

| Principle | Implementation |
| :--- | :--- |
| **Stateless API** | The Go layer does not store user sessions locally. All state is stored in Redis (refresh tokens, rate-limit counters) or encoded in the JWT itself, allowing horizontal scaling. |
| **Defense in Depth** | Security checks exist at Nginx (rate limits), Go (authorization middleware), PostgreSQL (RLS), and application query filters. |
| **Decoupled Logging** | Audit writes are pushed to a Redis queue and consumed by a background worker. This prevents slow disk I/O from blocking API responses. |
| **Database as Source of Truth** | Caches (Redis) are considered transient and can be rebuilt from the database at any time. No business-critical state lives exclusively in Redis. |
| **Configuration over Code** | Connections, timeouts, and rate limits are externally configurable via environment variables, avoiding hard-coded values. |

---

## Scalability & High Availability

| Component | Scalability Strategy |
| :--- | :--- |
| **Go API** | Horizontally scalable behind Nginx. Instances are fully stateless and can be added/removed dynamically. |
| **Next.js** | Horizontally scalable; static assets are offloaded to Nginx to reduce server load. |
| **PostgreSQL** | Supports read replicas for analytical queries; writes remain on the primary. Connection pooling is enabled. |
| **Redis** | Can be clustered for high availability; supports persistence for recovery. |

---

## Infrastructure Security (Perimeter)

| Layer | Measure |
| :--- | :--- |
| **Network** | TLS enforced for all traffic; HSTS header ensures HTTPS-only connections. |
| **API Gateway (Nginx)** | DDoS protection, request size limits, and IP-based rate limiting. |
| **Application** | Strict CORS allowlist to prevent cross-origin attacks. |
| **Data Storage** | Backups are encrypted at rest. Database credentials are stored in secure environment variables, not in code. |

---

## References

This document complements the architectural specifications:

- [auth.md](./auth.md) – Authentication and authorization overview.
- [authentication-process.md](./authentication-process.md) – Login, tokens, sessions.
- [database-logging.md](./database-logging.md) – Audit trail implementation.
