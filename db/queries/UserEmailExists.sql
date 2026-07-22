SELECT EXISTS email
FROM erp.users
WHERE email = $1;
