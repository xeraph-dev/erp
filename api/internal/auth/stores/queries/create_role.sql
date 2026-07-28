INSERT INTO auth.roles (role_name, role_level)
VALUES (@role_name, @role_level)
RETURNING id, role_name, role_level;
