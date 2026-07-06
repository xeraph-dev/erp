# Users (Identity Registry)

The `users` entity serves as the central identity registry for the system. Every authenticated interaction traces back to a user record, making it the foundation for authentication, authorization, and auditability.

---

## Core Concept

The user record represents an individual actor in the system. It carries:

- **Identity information** (unique username, email).
- **Credentials** (password hash).
- **Authorization context** (role).
- **Lifecycle state** (active, soft‑deleted).
- **Audit trail** (creation and modification timestamps, last login).

All other system entities (orders, logs, tokens, etc.) reference the user for accountability and ownership.

---

## User Lifecycle

| Phase              | Description                                                                                                                                                       | Responsible Component                                                                        |
| :----------------- | :---------------------------------------------------------------------------------------------------------------------------------------------------------------- | :------------------------------------------------------------------------------------------- |
| **Registration**   | A new user is created – either via self‑registration (using an invitation token) or by an administrator. The password is hashed and the default role is assigned. | [password-security.md](./password-security.md), [invitation-token.md](./invitation-token.md) |
| **Authentication** | The user provides credentials; the system verifies them and issues a session.                                                                                     | [authentication-process.md](./authentication-process.md)                                     |
| **Authorization**  | Every request checks the user's role/permissions.                                                                                                                 | [roles-and-permissions.md](./roles-and-permissions.md)                                       |
| **Modification**   | Users (or admins) can update non‑sensitive fields. Passwords have a dedicated change flow.                                                                        | [password-security.md](./password-security.md)                                               |
| **Deletion**       | Accounts are soft‑deleted, then optionally restored or purged.                                                                                                    | [soft-delete.md](./soft-delete.md)                                                           |

---

## Ownership and Data Scoping

- **By default**, users can access and modify only their own data (own profile, own orders, etc.).
- **Administrative users** can access and manage all data, subject to explicit permission checks.
- Data isolation is enforced at the query level (application and database layers) to ensure that a user never sees another user's data unless authorized.

---

## Administrative Capabilities

Users with appropriate permissions (typically `admin` or roles with `users:*` permissions) can:

- List, view, create, and update any user.
- Assign or change a user's role.
- Soft‑delete and restore user accounts.
- Initiate password resets on behalf of users.

Key constraints:

- No one can delete the last active `admin` user.
- No user can assign themselves a higher role than their own.
- A user cannot soft‑delete themselves.

---

## Self‑Service Capabilities

Authenticated users can manage their own profile within boundaries:

- View and update their own username, email, and other non‑sensitive fields.
- Change their own password (requires current password for verification).
- Delete their own account (soft‑delete, with confirmation).

---

## Data Integrity and Uniqueness

- Usernames and emails must be globally unique **among active records**.
- Soft‑deleted records do not block reuse of their username or email, allowing future users to take them.
- All foreign keys referencing users (e.g., logs, tokens) remain intact even after soft‑deletion, preserving audit chains.

---

## Password Management

Password‑related responsibilities are fully delegated to the password security module:

| Responsibility                             | Reference                                                |
| :----------------------------------------- | :------------------------------------------------------- |
| Hashing on creation/change                 | [password-security.md](./password-security.md)           |
| Verification during login                  | [password-security.md](./password-security.md)           |
| Strength validation                        | [password-security.md](./password-security.md)           |
| Reset flow                                 | [password-security.md](./password-security.md)           |
| Session invalidation after password change | [authentication-process.md](./authentication-process.md) |

---

## Audit Integration

All user‑related security events are logged:

- Account creation (via registration or admin).
- Profile changes (username, email).
- Role assignments.
- Password changes and resets.
- Account deletion and restoration.

For details, see [database-logging.md](./database-logging.md).

---

## References

- [authentication-process.md](./authentication-process.md) – Login and session issuance.
- [password-security.md](./password-security.md) – Password handling.
- [roles-and-permissions.md](./roles-and-permissions.md) – Role assignments and permission checks.
- [soft-delete.md](./soft-delete.md) – Account deletion/restoration lifecycle.
- [invitation-token.md](./invitation-token.md) – Registration gatekeeping via tokens.
- [database-logging.md](./database-logging.md) – Auditing all user actions.
- [modification-fields.md](./modification-fields.md) – Creation and modification timestamps.
