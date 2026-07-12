INSERT INTO erp.users (name, password, email)
VALUES ($1, $2, $3)
RETURNING *;
