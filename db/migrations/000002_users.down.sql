DROP TRIGGER IF EXISTS prevent_default_users_modification ON erp.users_table;
DROP FUNCTION IF EXISTS erp.prevent_default_users_modification;


DROP VIEW IF EXISTS erp.users;


DROP INDEX IF EXISTS active_users_email_idx;
DROP INDEX IF EXISTS active_users_name_idx;


DROP TRIGGER IF EXISTS default_deleted_at ON erp.users_table;
DROP FUNCTION IF EXISTS erp.default_deleted_at;


ALTER TABLE erp.users_table
DROP COLUMN IF EXISTS created_by_id,
DROP COLUMN IF EXISTS updated_by_id,
DROP COLUMN IF EXISTS deleted_by_id;


ALTER TABLE erp.base_table
DROP COLUMN IF EXISTS created_by_id,
DROP COLUMN IF EXISTS updated_by_id,
DROP COLUMN IF EXISTS deleted_by_id;


DROP FUNCTION IF EXISTS erp.system_user_id;


DROP MATERIALIZED VIEW IF EXISTS erp.users_system;


DELETE FROM erp.users_table
WHERE name IN ('system', 'admin');


DROP TRIGGER IF EXISTS update_updated_at ON erp.users_table;
DROP FUNCTION IF EXISTS erp.update_updated_at;


DROP TABLE IF EXISTS erp.users_table;
