SELECT EXISTS (
    SELECT 1
    FROM erp.users
    WHERE username = $1
);
