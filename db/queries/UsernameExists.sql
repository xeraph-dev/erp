SELECT EXISTS username
FROM erp.users
WHERE username = $1;
