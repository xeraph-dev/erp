DROP TRIGGER IF EXISTS prevent_default_roles_modification ON erp.roles_table;
DROP FUNCTION IF EXISTS erp.prevent_default_roles_modification;


DELETE FROM erp.roles_table
WHERE name IN ('admin', 'user');


DROP VIEW IF EXISTS erp.roles;


DROP INDEX IF EXISTS active_users_roles_idx ON erp.roles_table;


DROP TRIGGER IF EXISTS default_deleted_at ON erp.roles_table;


DROP TRIGGER IF EXISTS update_updated_at ON erp.roles_table;


DROP TABLE IF EXISTS erp.roles_table;
