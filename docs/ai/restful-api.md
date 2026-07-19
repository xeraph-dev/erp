# RESTful API Design (Architectural Standards)

This document defines the high-level design principles and standards for all public and internal APIs. It ensures consistency across endpoints, predictability for clients, and alignment with the broader system architecture.

---

## Core Principles

- **Resource‑Oriented:** APIs are organized around **resources** (nouns), not actions (verbs). Actions are expressed via HTTP methods.
- **Stateless:** Each request contains all necessary information. The server does not maintain client state between requests (except for session cookies handled by authentication).
- **JSON Payloads:** All request bodies and responses use `application/json`. XML and other formats are not supported.
- **Human‑Readable & Machine‑Friendly:** URL structures are intuitive for developers while maintaining strict machine parsability.

---

## URL Structure

- **Base Path:** All API endpoints are prefixed with `/api/v{version}/` (e.g., `/api/v1/users`).
- **Resources:** Use plural, lowercase nouns (e.g., `/users`, `/orders`, `/products`).
- **Sub‑Resources:** Use nested paths for related entities (e.g., `/users/{id}/orders`).
- **Actions:** Avoid verbs in URLs. For operations that do not map to standard CRUD, use a dedicated sub‑resource with a verb (e.g., `POST /users/{id}/restore` for soft‑delete restoration) or a query parameter (`?action=...`) only as a last resort.

---

## HTTP Methods & Semantics

| Method | Semantic | Idempotent? | Safe? | Typical Use |
| :--- | :--- | :--- | :--- | :--- |
| `GET` | Retrieve a resource or collection. | Yes | Yes | Fetching data. |
| `POST` | Create a new resource, or execute a custom action. | No | No | Creating orders, users, or triggering actions like `restore`. |
| `PUT` | Replace an entire resource (full update). | Yes | No | Complete replacement of a record. |
| `PATCH` | Partial update of a resource. | No (can be) | No | Updating specific fields (e.g., user email). |
| `DELETE` | Remove a resource (maps to [soft-delete.md](./soft-delete.md)). | Yes | No | Deleting a record (logical deletion). |

**Idempotency Promise:**
- `GET`, `PUT`, and `DELETE` are strictly idempotent. Multiple identical requests produce the same result as a single request.
- `POST` is **not** idempotent by default, except where explicitly documented (e.g., custom idempotency keys for financial transactions).

---

## HTTP Status Codes

The system uses a standard, minimal set of status codes to indicate success, client error, and server error.

| Code | Meaning | Use Case |
| :--- | :--- | :--- |
| `200 OK` | Success | Standard success response with a body. |
| `201 Created` | Resource created | Successful `POST` that creates a new resource. |
| `204 No Content` | Success, no body | Successful `DELETE`, `PUT`, or `PATCH` where no content needs to be returned. |
| `400 Bad Request` | Malformed request | Validation errors, malformed JSON, missing required fields. |
| `401 Unauthorized` | Missing or invalid credentials | [authentication-process.md](./authentication-process.md) – JWT missing, expired, or invalid. |
| `403 Forbidden` | Authenticated but not authorized | [roles-and-permissions.md](./roles-and-permissions.md) – Insufficient permissions for the resource/action. |
| `404 Not Found` | Resource not found | The requested URI does not map to an existing resource. |
| `409 Conflict` | Data conflict | Unique constraint violation (e.g., duplicate username). |
| `422 Unprocessable Entity` | Semantic error | Business logic violation (e.g., trying to delete a user with active orders). |
| `429 Too Many Requests` | Rate‑limited | Exceeded rate limits (enforced by Nginx or application). |
| `500 Internal Server Error` | Server failure | Unexpected system error; client should retry after a delay. |

---

## Request & Response Format

### Request Headers

| Header | Requirement | Description |
| :--- | :--- | :--- |
| `Content-Type: application/json` | Required for `POST`, `PUT`, `PATCH` | Indicates the request body format. |
| `Accept: application/json` | Recommended | Indicates the desired response format (JSON). |
| `Cookie` | Required for authenticated endpoints | Carries the JWT access token (see [authentication-process.md](./authentication-process.md)). |

### Response Structure (Success)

- **Collections:** Include `data` (array of resources), `pagination` metadata (if applicable), and `total` count.
- **Single Resources:** Return the full resource object directly as the top‑level JSON.
- **No Content:** `204` responses have an empty body.

### Error Response Structure

All error responses (4xx and 5xx) follow a consistent envelope to simplify client error handling:

- A top‑level `error` object containing:
  - `code`: A machine‑readable string (e.g., `validation_failed`, `permission_denied`).
  - `message`: A human‑readable summary.
  - `details` (optional): Additional context (e.g., specific field validation failures, nested errors).

---

## Versioning Strategy

- **URL Versioning:** The major version is embedded in the path (e.g., `/api/v1/users`).
- **Backward Compatibility:** Within the same major version (`v1`), all changes are **additive only**:
  - New fields can be added to responses.
  - New endpoints can be introduced.
  - Existing fields and endpoints **cannot** be removed or have their semantics changed.
- **Breaking Changes:** Necessitate a new major version (`v2`). Both versions can coexist simultaneously during migration.
- **Deprecation Policy:** Deprecated endpoints are marked with a `Deprecation` header and remain available for at least one full release cycle before removal.

---

## Pagination, Filtering & Sorting

For collection endpoints (`GET /api/v1/resources`), query parameters standardize data retrieval.

| Parameter | Example | Description |
| :--- | :--- | :--- |
| `page` | `?page=2` | Page number (1‑based). |
| `limit` | `?limit=50` | Number of items per page. A default and a maximum are enforced. |
| `sort` | `?sort=-created_at` | Sort order: `+` for ascending, `-` for descending. |
| `filter` | `?filter[status]=active` | Filtering by specific fields. The exact filterable fields are documented per resource. |
| `include` | `?include=user,role` | Eager loading of related resources (to reduce round trips). |

**Defaults:** If `page` and `limit` are omitted, a sensible default (e.g., `page=1`, `limit=25`) is applied.

---

## Integration with System Features

| Feature | API Integration |
| :--- | :--- |
| [authentication-process.md](./authentication-process.md) | Authentication is handled via HTTP‑only cookies, not `Authorization` headers. The API reads the JWT from the `Cookie` header. |
| [roles-and-permissions.md](./roles-and-permissions.md) | Authorization middleware inspects the user's role/permissions before routing. `403` responses are returned for missing permissions. |
| [soft-delete.md](./soft-delete.md) | `DELETE` requests perform logical deletion. `POST /{resource}/{id}/restore` is the dedicated endpoint for restoration. |
| [database-logging.md](./database-logging.md) | All `POST`, `PUT`, `PATCH`, and `DELETE` operations are automatically audited. |
| [users.md](./users.md) | The user resource serves as the identity anchor for all ownership and scoping rules. |

---

## Security & Compliance

- **CORS:** Strict origin allowlist. Only configured frontend domains can call the API.
- **Rate Limiting:** Applied at the Nginx layer for brute‑force protection (login, registration, password reset) and at the application layer for resource‑intensive endpoints.
- **Sensitive Data:** Fields like `password_hash` are **never** exposed in API responses.
- **Idempotency Keys (Optional):** For financial or critical operations (e.g., payment creation), the API supports an `Idempotency-Key` header to prevent duplicate processing.

---

## References

- [auth.md](./auth.md) – Overall authentication/authorization architecture.
- [authentication-process.md](./authentication-process.md) – JWT and cookie handling specifics.
- [users.md](./users.md) – The primary resource entity.
- [soft-delete.md](./soft-delete.md) – Semantics of `DELETE` and restoration.
- [database-logging.md](./database-logging.md) – Audit logging for API requests.
- [roles-and-permissions.md](./roles-and-permissions.md) – Permission mapping and enforcement.
