DROP TRIGGER IF EXISTS soft_delete_roles ON auth.roles;
DROP FUNCTION IF EXISTS auth.soft_delete_roles;


DROP VIEW IF EXISTS auth.roles;


DROP TRIGGER IF EXISTS update_modification_fields ON auth.__roles;


DROP FUNCTION IF EXISTS auth.user_role_id;
DROP FUNCTION IF EXISTS auth.admin_role_id;


DELETE FROM auth.__roles
WHERE role_name IN ('admin', 'user');


DROP TABLE IF EXISTS auth.__roles;
