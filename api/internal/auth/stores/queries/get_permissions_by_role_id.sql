SELECT id, permission_name
FROM auth.permissions
INNER JOIN auth.permissions_roles ON permissions.id = permissions_roles.permission_id
WHERE permissions_roles.role_id = @role_id;
