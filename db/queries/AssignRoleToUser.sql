INSERT INTO erp.roles_users (role_id, user_id)
VALUES ($1, $2)
RETURNING role_id, user_id;
