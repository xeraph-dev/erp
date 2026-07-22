SELECT EXISTS (
    SELECT 1
    FROM erp.users
    WHERE email = $1
);
