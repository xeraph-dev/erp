# Modification Fields (Timestamps)

All data tables in the system automatically track when each record was created and when it was last modified. This provides a consistent, auditable timeline for every piece of data.

---

## Core Concept

Every table **must** include:

- **`created_at`** – the exact moment the record was inserted. This value is set once and never changes.
- **`updated_at`** – the exact moment the record was last modified. This value is automatically refreshed on every data change.

These fields are maintained at the **database level** to ensure consistency, regardless of the application layer or access method.

---

## Field Behavior

| Field            | Behavior                                                                                                                  |
| :--------------- | :------------------------------------------------------------------------------------------------------------------------ |
| **`created_at`** | Set automatically on insertion. Immutable after creation – any attempt to modify it is rejected.                          |
| **`updated_at`** | Set automatically on insertion. Automatically updated to the current time whenever any other field in the record changes. |

---

## Modification Tracking Principles

- **Automatic:** The database handles both fields without requiring explicit application logic.
- **Immutable Creation:** `created_at` is permanently frozen after the record is created.
- **State‑Change Driven:** `updated_at` is bumped only when the record's actual data changes. No‑op updates (changing a field to its current value) do not trigger an update.
- **Comprehensive:** Both fields are included in all snapshots for audit and forensics, providing a complete timeline of the record's lifecycle.

---

## Integration with Other Features

| Feature                                      | Interaction                                                                                                                                                              |
| :------------------------------------------- | :----------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [soft-delete.md](./soft-delete.md)           | Soft‑deleting a record (setting `deleted_at`) also updates `updated_at` to reflect the state change. Restoring a record similarly updates `updated_at`.                  |
| [database-logging.md](./database-logging.md) | Both timestamps are included in the before/after snapshots for every `UPDATE` operation. Changes to `updated_at` are automatically captured as part of any modification. |
| [users.md](./users.md)                       | These fields are present in the user record, tracking account creation and profile modifications.                                                                        |
| API Contracts                                | These fields are read‑only in all API responses. Incoming requests must not attempt to supply values for them.                                                           |

---

## Querying and Indexing

- **Default visibility:** Both fields are returned in standard API responses.
- **Time‑based filtering:** The system supports filtering by these fields (e.g., records created after a certain date, modifications since yesterday).
- **Performance:** Appropriate database indexing is applied to support common time‑based query patterns without degrading performance.

---

## Invariants

| Invariant                                             | Enforcement                                                 |
| :---------------------------------------------------- | :---------------------------------------------------------- |
| `created_at` is never null.                           | Database-level default.                                     |
| `updated_at` is never null.                           | Database-level default.                                     |
| `created_at` never changes after insertion.           | Database-level rejection of any modification attempt.       |
| `updated_at` changes on every real data modification. | Database-level automatic refresh.                           |
| `updated_at` does not change on no‑op updates.        | Database-level change detection prevents unnecessary bumps. |

---

## References

- [users.md](./users.md) – The user table incorporates these fields as part of the core schema.
- [soft-delete.md](./soft-delete.md) – The `deleted_at` field coexists with these modification fields but is managed separately.
- [database-logging.md](./database-logging.md) – Both fields are included in audit snapshots for complete historical reconstruction.
