SELECT id, username, email password_hash, first_name, last_name
FROM auth.users
WHERE email = @email;
