INSERT INTO auth.users (username, email, password_hash, first_name, last_name)
VALUES (@username, @email, @password_hash, @first_name, @last_name)
RETURNING id, username, email, password_hash, first_name, last_name;
