SELECT id, role_name, role_level
FROM auth.roles
WHERE id = @id;
