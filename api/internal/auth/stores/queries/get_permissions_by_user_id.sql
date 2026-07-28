SELECT id, permission_name
FROM auth.permissions
INNER JOIN auth.permissions_users ON permissions.id = permissions_users.permission_id
WHERE permissions_users.user_id = @user_id;
