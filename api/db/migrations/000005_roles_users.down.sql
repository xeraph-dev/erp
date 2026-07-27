DELETE FROM auth.roles_users
WHERE role_id = auth.admin_role_id()
  AND user_id = (SELECT id FROM auth.users WHERE username = 'admin' LIMIT 1);


DELETE FROM auth.__users
WHERE username = 'admin';


DROP TRIGGER IF EXISTS update_modification_fields ON auth.roles;


DROP TABLE IF EXISTS auth.roles_users;
