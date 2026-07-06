# Authentication Process (Identity Verification)

This document describes how the system verifies a user's identity and establishes a secure session. It covers login, logout, token issuance, cookie handling, session refresh, and invalidation.

---

## Core Concept

Authentication answers: **"Who is the user?"**

The system uses a **stateless, token‑based** authentication model:

- Identity is verified once (during login), and the result is encoded into a self‑contained token.
- The server does not store session state (except for refresh tokens).
- All subsequent requests carry the token, which the server validates independently.

---

## Authentication Flow (Login)

1. The user submits credentials (username/email + password) via a secure endpoint.
2. The system validates the credentials against the stored password hash.
3. Upon success:
   - The system issues two tokens:
     - A **short‑lived access token** (for API requests).
     - A **longer‑lived refresh token** (for obtaining new access tokens).
   - Both tokens are placed in **secure, HTTP‑only cookies**.
   - The user's last login timestamp is updated.
4. The user is now authenticated for subsequent requests.

---

## Token Strategy

| Token Type        | Purpose                                                      | Characteristics                                                                         |
| :---------------- | :----------------------------------------------------------- | :-------------------------------------------------------------------------------------- |
| **Access Token**  | Used for every authenticated API request.                    | Short‑lived (minutes to hours); self‑contained (contains user identity and claims).     |
| **Refresh Token** | Used to obtain a new access token without re‑authenticating. | Longer‑lived (hours to days); stored server‑side for invalidation; rotated on each use. |

This dual‑token approach balances security (short exposure of access tokens) and usability (seamless session extension).

---

## Cookie Transmission

- **All tokens** are transmitted via **HTTP‑only, Secure, SameSite** cookies.
- This prevents:
  - Cross‑site scripting (XSS) – tokens are inaccessible to JavaScript.
  - Cross‑site request forgery (CSRF) – SameSite attribute restricts cross‑origin requests.
- Cookies are set with appropriate path and domain restrictions.

---

## Session Refresh

- When the access token expires, the client can request a new one using the refresh token.
- The refresh token is validated, and a new access token (and optionally a new refresh token) is issued.
- Refresh tokens are **rotated** on each use – the old refresh token is invalidated and a new one is issued.
- This reduces the window for replay attacks.

---

## Logout

- The client signals logout.
- The system **clears the token cookies** (sets expiration to a past date).
- The refresh token is **invalidated server‑side** to prevent reuse.
- The session is effectively terminated.

---

## Session Invalidation

All active sessions are invalidated in the following scenarios:

- **Password change or reset** – the user must re‑authenticate with the new password.
- **Logout** – explicit termination.
- **Account soft‑deletion** – the user cannot log in until restored.
- **Admin‑initiated revocation** – administrators can force a logout for specific users.

Invalidation mechanisms include:

- Clearing cookies (client‑side).
- Blacklisting or deleting refresh tokens (server‑side).
- (Access tokens are short‑lived and are not revoked; they are allowed to expire naturally.)

---

## Authentication Middleware

Every protected request passes through an authentication middleware that:

1. Extracts the access token from the cookie.
2. Validates the token's signature, expiration, and integrity.
3. Injects the authenticated `user_id` (and optionally permissions) into the request context.
4. Rejects the request (401 Unauthorized) if validation fails.

The middleware does **not** perform authorization (permission checks) – that is handled separately by the authorization middleware.

---

## Integration with Other Features

| Feature                                                | Integration Point                                                         |
| :----------------------------------------------------- | :------------------------------------------------------------------------ |
| [password-security.md](./password-security.md)         | Credential verification and password hashing.                             |
| [users.md](./users.md)                                 | User identity lookup and status (active/deleted).                         |
| [roles-and-permissions.md](./roles-and-permissions.md) | The authenticated user's permissions are cached for authorization checks. |
| [database-logging.md](./database-logging.md)           | Login attempts (success/failure) and logout events are audited.           |
| [soft-delete.md](./soft-delete.md)                     | Soft‑deleted users are denied authentication.                             |
| [modification-fields.md](./modification-fields.md)     | The user's `last_login_at` is updated on each successful login.           |

---

## Security Properties

| Property                        | Implementation Principle                                                                                     |
| :------------------------------ | :----------------------------------------------------------------------------------------------------------- |
| **Confidentiality**             | Tokens are transmitted only over TLS, never in URLs or plaintext.                                            |
| **Integrity**                   | Tokens are digitally signed to prevent tampering.                                                            |
| **Replay Resistance**           | Short lifetimes and refresh token rotation limit replay windows.                                             |
| **Session Fixation Prevention** | New tokens are issued on login; no session identifier is reused.                                             |
| **Stateless Scalability**       | The access token is self‑contained, eliminating server‑side session storage and enabling horizontal scaling. |

---

## References

- [password-security.md](./password-security.md) – Password hashing and verification.
- [users.md](./users.md) – User identity and account state.
- [roles-and-permissions.md](./roles-and-permissions.md) – Authorization after authentication.
- [database-logging.md](./database-logging.md) – Audit of authentication events.
- [invitation-token.md](./invitation-token.md) – Pre‑registration flow (not authentication).
- [modification-fields.md](./modification-fields.md) – Timestamp for last login.
