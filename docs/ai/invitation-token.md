# Invitation Token (Registration Gatekeeping)

The system restricts new user registration to explicitly invited individuals. This prevents automated bot sign-ups and ensures that only authorized parties can create accounts.

---

## Core Concept

Registration requires a valid, single-use invitation token. Tokens are:

- **Generated** exclusively by administrators.
- **Time-bound** – they expire after a configurable period.
- **Single-use** – consumed immediately upon successful registration.
- **Traceable** – linked to the admin who created them and the user who used them.

---

## Token Lifecycle

| Phase            | Description                                                                                                                                        | Responsible Owner                                                          |
| :--------------- | :------------------------------------------------------------------------------------------------------------------------------------------------- | :------------------------------------------------------------------------- |
| **Generation**   | An administrator requests a new token via a privileged endpoint. The system creates a cryptographically random token with a fixed expiration date. | [roles-and-permissions.md](./roles-and-permissions.md) (admin permissions) |
| **Distribution** | The administrator receives the token string and shares it with the intended user through a secure out-of-band channel (e.g., email, secure link).  | Administrator                                                              |
| **Validation**   | During registration, the user submits the token. The system checks that the token exists, has not been used, and has not expired.                  | Registration flow                                                          |
| **Consumption**  | Upon successful user creation, the system permanently marks the token as used, preventing any further reuse.                                       | Transactional registration                                                 |
| **Expiration**   | Unused tokens automatically become invalid after their expiration time. Expired tokens are rejected during validation.                             | System-defined expiration policy                                           |

---

## Integration with Registration

The invitation token is a mandatory prerequisite for self-registration. The flow is:

1. The user submits the token alongside their chosen credentials.
2. The system validates the token in a single atomic step with the user creation process.
3. If validation fails (invalid, expired, or already used), the registration is rejected.
4. If validation passes, the user account is created and the token is marked as used within the same database transaction.

This ensures that a token cannot be consumed if user creation fails, nor can a user be created without a valid token.

---

## Optional Restrictions

The system optionally supports restricting a token to a specific email address. If set, the registration flow will only accept that exact email, adding an extra layer of control over who can use the token.

---

## Auditability and Cleanup

- **Audit Trail:** Every token generation and consumption is recorded in the central audit log. This provides accountability for who invited whom and when.
- **Background Purge:** Expired and consumed tokens are periodically removed from the active registry to prevent unbounded growth, while maintaining historical records for compliance where required.

---

## Security Properties

| Property                   | Implementation Principle                                                                                             |
| :------------------------- | :------------------------------------------------------------------------------------------------------------------- |
| **Unguessability**         | Tokens are generated using a cryptographically secure random source, making them infeasible to guess or brute-force. |
| **Single-Use**             | Tokens are atomically consumed during registration, preventing replay attacks.                                       |
| **Time-Bounded**           | Automatic expiration limits the window of opportunity for misuse.                                                    |
| **Enumeration Resistance** | The validation endpoint returns generic error messages, preventing attackers from probing for valid tokens.          |

---

## References

- [authentication-process.md](./authentication-process.md) – The registration flow where the token is consumed.
- [roles-and-permissions.md](./roles-and-permissions.md) – Defines the `tokens:generate` permission required for token creation.
- [database-logging.md](./database-logging.md) – Logs token generation and consumption events.
- [users.md](./users.md) – The user record created upon successful token consumption.
- [password-security.md](./password-security.md) – Credential validation that occurs alongside token validation during registration.
