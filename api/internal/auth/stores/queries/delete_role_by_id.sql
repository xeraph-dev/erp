DELETE FROM auth.roles
WHERE id = @id
RETURNING id, role_name, role_level;
