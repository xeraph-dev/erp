DROP TRIGGER IF EXISTS soft_delete_roles ON erp.roles;
DROP FUNCTION IF EXISTS erp.soft_delete_roles;


DROP VIEW IF EXISTS erp.roles;


DROP TRIGGER IF EXISTS update_modification_fields ON erp.roles_table;


DROP INDEX IF EXISTS active_roles_role_name_idx;


DELETE FROM erp.roles_table
WHERE role_name IN ('admin', 'user');


DROP TABLE IF EXISTS erp.roles_table;
