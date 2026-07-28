DELETE FROM auth.roles_users
WHERE role_id = auth.admin_role_id()
  AND user_id = (SELECT id FROM auth.users WHERE username = 'admin' LIMIT 1);


DELETE FROM auth.__users
WHERE username = 'admin';


DROP TRIGGER IF EXISTS update_modification_fields ON auth.roles;


DROP TABLE IF EXISTS auth.roles_users;


DROP TRIGGER IF EXISTS soft_delete_roles ON auth.roles;
DROP FUNCTION IF EXISTS auth.soft_delete_roles;


DROP VIEW IF EXISTS auth.roles;


DROP TRIGGER IF EXISTS update_modification_fields ON auth.__roles;


DROP FUNCTION IF EXISTS auth.user_role_id;
DROP FUNCTION IF EXISTS auth.admin_role_id;


DELETE FROM auth.__roles
WHERE role_name IN ('admin', 'user');


DROP TABLE IF EXISTS auth.__roles;
