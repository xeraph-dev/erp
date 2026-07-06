# Soft Delete (Logical Deletion)

The system implements soft deletion for critical data to prevent accidental data loss, maintain referential integrity, and preserve audit history. Data is logically removed but physically retained for recovery and compliance purposes.

---

## Core Concept

Deleting a record does not physically remove it from the database. Instead, the record is marked as inactive:

- **Active records** are visible and fully functional.
- **Deleted records** are hidden from all standard operations but remain stored.

This approach allows for:

- **Accidental deletion recovery** – administrators can restore deleted records.
- **Audit preservation** – the record's history remains intact.
- **Unique value reuse** – usernames, emails, and other unique identifiers can be reused by new records after deletion.

---

## Deletion Behavior

The standard `DELETE` operation always performs a **soft delete**:

- The record is marked with a deletion timestamp and the identifier of the user who performed the deletion.
- The record is excluded from all normal queries and operations.
- The operation is idempotent – attempting to delete an already-deleted record succeeds without further changes.
- **Hard (physical) deletion** is never triggered by the standard API.

---

## Restoration

Soft-deleted records can be **restored** via a dedicated privileged operation:

- Restoration clears the deletion markers, making the record active again.
- Restoration is restricted to users with appropriate administrative permissions.
- The restoration event is logged for accountability.

---

## Physical Removal (Purge)

Permanent physical deletion is a **separate, privileged operation**:

- It is **not** exposed through standard APIs.
- Physical removal is allowed only through:
  - Dedicated administrative tools or internal endpoints.
  - Compliance-driven data erasure requests (e.g., GDPR "right to be forgotten").
- Purge operations are subject to strict controls, typically requiring additional authentication (e.g., multi-factor) and are logged as critical security events.

---

## Querying Rules

- **Default:** All standard queries automatically exclude soft-deleted records. Applications see only active data.
- **Opt‑in:** Administrators can explicitly request to include deleted records when needed (e.g., for investigation or restoration).
- This behavior is enforced at the data access layer to prevent accidental exposure of deleted records.

---

## Data Integrity Constraints

- **Unique constraints** (e.g., usernames, emails) are designed to **exclude soft-deleted records**. This allows a new account to use a value previously held by a deleted account, avoiding unnecessary friction.
- All application-level logic that checks uniqueness is scoped to active records only.

---

## Audit Integration

All soft‑delete and restoration operations are logged in the audit trail, capturing:

- The user performing the deletion or restoration.
- The target record.
- The timestamp.
- The state before and after the operation.

For details, see [database-logging.md](./database-logging.md).

---

## References

- [users.md](./users.md) – The user entity includes soft-delete fields and supports deletion/restoration.
- [database-logging.md](./database-logging.md) – Auditing of delete and restore operations.
- [modification-fields.md](./modification-fields.md) – The `deleted_at` field coexists with `created_at` and `updated_at`.
- [roles-and-permissions.md](./roles-and-permissions.md) – Defines the `users:delete` and `users:restore` permissions.
- [authentication-process.md](./authentication-process.md) – Deleted users cannot authenticate (blocked at login).
