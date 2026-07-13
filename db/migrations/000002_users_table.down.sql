DROP TRIGGER IF EXISTS update_modification_fields ON erp.users_table;


DROP FUNCTION IF EXISTS erp.update_modification_fields;


DROP INDEX IF EXISTS active_users_email_idx;
DROP INDEX IF EXISTS active_users_username_idx;


ALTER TABLE erp.users_table
DROP COLUMN IF EXISTS created_by_id,
DROP COLUMN IF EXISTS updated_by_id,
DROP COLUMN IF EXISTS deleted_by_id;


DROP FUNCTION IF EXISTS erp.system_user_id();


DELETE FROM erp.users_table
WHERE username = 'system';

DROP TABLE IF EXISTS erp.users_table;
