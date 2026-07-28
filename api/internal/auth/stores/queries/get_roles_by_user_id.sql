SELECT roles.id, roles.role_name, roles.role_level
FROM auth.roles
INNER JOIN auth.roles_users ON roles.id = roles_users.role_id
WHERE roles_users.user_id = @user_id;
