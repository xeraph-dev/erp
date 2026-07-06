# Authentication & Authorization (Architectural Overview)

This document outlines the high-level architectural separation between **identity verification** (Authentication) and **access control** (Authorization). All implementation specifics are delegated to the referenced sub-specifications.

---

## 1. Component References

| Document                                                 | Scope                                                                         |
| :------------------------------------------------------- | :---------------------------------------------------------------------------- |
| [password-security.md](./password-security.md)           | Credential hashing, password policies, and reset/change workflows.            |
| [invitation-token.md](./invitation-token.md)             | Single-use registration tokens, generation, and consumption rules.            |
| [authentication-process.md](./authentication-process.md) | Login, logout, JWT issuance, cookie handling, and token refresh logic.        |
| [roles-and-permissions.md](./roles-and-permissions.md)   | RBAC model, role hierarchies, permission definitions, and enforcement layers. |
| [database-logging.md](./database-logging.md)             | Audit trail for security-critical events.                                     |
| [soft-delete.md](./soft-delete.md)                       | User lifecycle termination (logical vs. physical removal).                    |
| [users.md](./users.md)                                   | The user entity as the central identity registry.                             |

---

## 2. Complete User Lifecycle

All system interactions follow a controlled sequence from invitation to account termination:

| Phase              | Action                                                                           | Owner                                                                                         |
| :----------------- | :------------------------------------------------------------------------------- | :-------------------------------------------------------------------------------------------- |
| **Invitation**     | Admin generates a cryptographically strong, time-bound token.                    | [invitation-token.md](./invitation-token.md)                                                  |
| **Registration**   | New user submits the token along with credentials to activate the account.       | [invitation-token.md](./invitation-token.md) + [password-security.md](./password-security.md) |
| **Authentication** | User provides credentials; system verifies them and issues a session identifier. | [authentication-process.md](./authentication-process.md)                                      |
| **Session**        | Identity is maintained via short-lived tokens with a refresh mechanism.          | [authentication-process.md](./authentication-process.md)                                      |
| **Authorization**  | Every request is checked against the user's assigned permissions.                | [roles-and-permissions.md](./roles-and-permissions.md)                                        |
| **Auditing**       | All security events are recorded immutably.                                      | [database-logging.md](./database-logging.md)                                                  |
| **Account End**    | Account is logically deactivated (soft-delete) with restoration capability.      | [soft-delete.md](./soft-delete.md)                                                            |

---

## 3. Authentication (Identity Verification)

Authentication answers: **"Who is the user?"**

- **Architecture:** Stateless by design. The server does not maintain local session state beyond the issued tokens.
- **Transport:** Credentials and tokens are transmitted exclusively via secure, HTTP-only cookies to mitigate XSS/CSRF risks.
- **Lifecycle:** The system uses a dual-token strategy (short-lived access tokens + longer-lived refresh tokens) to balance security and usability.
- **Recovery:** Password resets rely on time-limited, out-of-band tokens sent via verified channels (e.g., email).

---

## 4. Authorization (Access Control)

Authorization answers: **"What is the user allowed to do?"**

- **Model:** Role-Based Access Control (RBAC). Every user is assigned exactly one role; every role is a collection of discrete permissions.
- **Philosophy:** Permissions are strictly additive (allow-only). The absence of an explicit permission implies a denial.
- **Defense-in-Depth:** The system enforces authorization at **three independent layers** to prevent bypasses:

1. **API Middleware** – blocks unauthorized requests at the entry point.
2. **Database Row-Level Security (RLS)** – ensures queries only return allowed data, even if the middleware fails.
3. **Repository Filtering** – adds explicit data-scoping constraints as a final safeguard.

---

## 5. Infrastructure Security

The application layer is protected by standard infrastructure practices:

| Aspect             | Decision                                                                                    |
| :----------------- | :------------------------------------------------------------------------------------------ |
| **Network**        | All traffic is encrypted via TLS; HSTS is enforced.                                         |
| **Throttling**     | Critical endpoints (login, registration, password reset) are subject to strict rate limits. |
| **Origin Control** | CORS is restricted to a specific allowlist of trusted origins.                              |
| **Data at Rest**   | Backups containing credentials are encrypted before storage.                                |

---

## 6. Audit & Monitoring

Security is not complete without traceability.

- All **security-critical actions** (login attempts, role assignments, password changes, user deletions) are automatically sent to the centralized audit log.
- The audit trail is immutable (append-only) and designed to support forensic analysis and compliance reporting. For full details, refer to [database-logging.md](./database-logging.md).
