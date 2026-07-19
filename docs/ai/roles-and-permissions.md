# Roles and Permissions (Authorization Model)

The system uses **Role-Based Access Control (RBAC)** to govern what authenticated users are allowed to do. Authorization is strictly separated from authentication and follows an additive (allow‑only) permission model.

---

## Core Concepts

- **Role:** A named collection of permissions. Each user is assigned exactly one role.
- **Permission:** A discrete capability to perform a specific action on a specific resource (e.g., `users:read`, `orders:write`).
- **Additive Model:** Permissions are _allow‑only_. The presence of a permission grants the action; its absence implies denial. Negative/deny rules are not supported.

---

## Role Hierarchy and Immutability

- The system is initialized with **system roles** (e.g., `admin`, `user`) that are **immutable** – they cannot be deleted or renamed, though their display names may be updated.
- Custom roles can be created, modified, and deleted by administrators.
- The default role for new users is the least‑privileged role (`user`).

---

## Permission Definition

- Permissions are **immutable**. They are seeded during system initialization and cannot be modified or deleted at runtime by any user, including administrators.
- Each permission follows a consistent naming pattern: `{resource}:{action}`.
- The set of available permissions is fixed and covers all system resources and operations.

---

## Assignment and Management

Administrators with the appropriate permission can:

- Create, update, and delete **custom roles** (system roles are protected).
- Assign or change a user's role.
- Modify the set of permissions attached to a role.

Key constraints:

- **System roles** (`admin`, `user`) cannot be deleted.
- An administrator cannot remove their own `admin` role (prevents accidental lockout).
- Permissions themselves are **not modifiable** – administrators can only assign/unassign existing permissions.

---

## Enforcement Layers (Defense‑in‑Depth)

Authorization is enforced at **three independent layers** to prevent bypasses:

1. **API Middleware:**
   - Checks the required permission for the incoming request against the user's cached permission set.
   - Returns `403 Forbidden` immediately if the permission is missing.

2. **Database Row‑Level Security (RLS):**
   - Restricts data visibility at the query level (e.g., users see only their own data, unless they have administrative permissions).
   - Ensures that even if the middleware is bypassed, the database rejects unauthorized reads/writes.

3. **Repository Filtering:**
   - The application repository layer adds explicit data‑scoping filters (e.g., `AND user_id = $1`) as a final safeguard.
   - This acts as a defense‑in‑depth measure against misconfigurations in the upper layers.

---

## Permission Caching

- Permissions are **cached** to minimize database load during request processing.
- The cache is invalidated immediately whenever a user's role changes.
- Cache lifetime is limited; on cache miss, the system retrieves permissions from the database.

---

## Data Visibility Rules

- By default, users can only access data that belongs to them (their own user record, their own orders, etc.).
- Users with administrative permissions can access all data, subject to explicit permission checks.

---

## Audit Integration

All role and permission changes are logged in the central audit trail, including:

- Who changed the role.
- Which user was affected.
- What the old and new roles were.

For details, see [database-logging.md](./database-logging.md).

---

## References

- [authentication-process.md](./authentication-process.md) – Authentication injects the `user_id` used for permission lookups.
- [users.md](./users.md) – The user record stores the assigned role.
- [database-logging.md](./database-logging.md) – Logging of role changes.
- [soft-delete.md](./soft-delete.md) – Restoring a user does not change their role.
- [modification-fields.md](./modification-fields.md) – Timestamps for role and permission modifications.
