DROP VIEW IF EXISTS erp.roles;


DROP INDEX IF EXISTS active_roles_idx ON erp.roles_table;


DROP TRIGGER IF EXISTS default_deleted_at ON erp.roles_table;


DROP TRIGGER IF EXISTS update_updated_at ON erp.roles_table;


DROP TABLE IF EXISTS erp.roles_table;
