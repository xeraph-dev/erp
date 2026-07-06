# Database Logging (Audit Trail)

The system maintains a complete, immutable history of all security-critical and data-modifying operations. This audit trail enables forensic analysis, compliance reporting, and reconstruction of past states.

---

## Core Concept

Every `CREATE`, `UPDATE`, and `DELETE` operation on tracked tables is recorded in a centralized, append-only audit log. Each entry captures:

- **What** changed (the affected table and record).
- **Who** performed the change (the authenticated user).
- **When** the change occurred.
- **What** the state was before and after the change.

---

## Logging Responsibility

The audit log is written by the **application layer**, not by database triggers.

**Rationale:**

- The application has direct access to the authenticated user context (user ID, IP address, correlation ID).
- Logging is decoupled from the main transaction to avoid degrading API performance.
- The application can selectively exclude internal or batch operations from the audit trail.

---

## Logging Flow

1. The application receives a request and authenticates the user.
2. The repository layer performs the data modification within a database transaction.
3. After the transaction successfully commits, the application constructs a log entry containing:
   - The affected table and record identifier.
   - The operation type (`CREATE`, `UPDATE`, `DELETE`).
   - A snapshot of the record before and after the change.
   - The authenticated user ID and client IP address.
4. The log entry is sent to a **decoupled, asynchronous queue** to be written to the audit table separately.
5. A background worker consumes the queue and performs batch writes to the audit log.

This ensures that slow disk I/O or logging failures never block or impact the critical path of API requests.

---

## Log Content Rules

| Operation  | Captured State                                                         |
| :--------- | :--------------------------------------------------------------------- |
| **CREATE** | Snapshot of the new record only.                                       |
| **UPDATE** | Full snapshot of the record before the change and after the change.    |
| **DELETE** | Snapshot of the record before removal (including soft-delete markers). |

The system logs the **entire record state**, not only the changed fields, to simplify reconstruction and forensics.

---

## Retention and Lifecycle

Audit logs are subject to a defined retention policy:

- **Hot Storage:** Recent logs (e.g., the last 30 days) are kept in the primary audit table for fast operational queries.
- **Cold Archive:** Older logs are moved to compressed archive storage or exported to external long-term storage.
- **Purge:** Logs beyond the legally required retention period are permanently deleted by a scheduled background process.

The retention period is configurable per deployment to meet varying compliance requirements.

---

## Access Control

- **Read access** to the audit log is strictly limited to users with elevated privileges (e.g., administrators with a specific permission).
- **No API** allows modification or deletion of log entries – the audit trail is append-only by design.
- **Sensitive data** within snapshots (e.g., password hashes) is either redacted before logging or encrypted at rest, ensuring credentials are never exposed in clear text even in the audit store.

---

## Integration with Other Features

The audit log integrates seamlessly with all major subsystems:

| Feature                                                  | Integration Point                                                                                                                             |
| :------------------------------------------------------- | :-------------------------------------------------------------------------------------------------------------------------------------------- |
| [authentication-process.md](./authentication-process.md) | Provides the `user_id` and IP address for every log entry.                                                                                    |
| [soft-delete.md](./soft-delete.md)                       | `DELETE` operations on soft-deletable tables capture the state before marking as deleted; `RESTORE` operations are logged as `UPDATE` events. |
| [roles-and-permissions.md](./roles-and-permissions.md)   | Role assignments and permission changes are logged with before/after states.                                                                  |
| [password-security.md](./password-security.md)           | Password changes are logged as events (the hash itself is never exposed).                                                                     |
| [invitation-token.md](./invitation-token.md)             | Token generation and consumption are logged for accountability.                                                                               |

---

## Invariants

| Invariant                                             | Enforcement                                                                                           |
| :---------------------------------------------------- | :---------------------------------------------------------------------------------------------------- |
| Log entries are never modified after creation.        | The audit table has no update/delete operations exposed via API.                                      |
| Every security-critical action generates a log entry. | Application code explicitly enqueues logs for relevant operations.                                    |
| Logging does not block the main request flow.         | Asynchronous queue decouples logging from the transaction commit.                                     |
| The user identity is always captured.                 | Authentication middleware injects `user_id` into the request context before any repository operation. |
