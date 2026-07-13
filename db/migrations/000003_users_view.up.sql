CREATE OR REPLACE VIEW erp.users AS
SELECT
    id,
    username,
    password_hash,
    email,
    first_name,
    last_name
FROM erp.users_table
WHERE
    username <> 'system'
    AND deleted_at IS NULL
    AND deleted_by_id IS NULL
WITH CHECK OPTION; -- noqa: disable=PRS


CREATE OR REPLACE FUNCTION erp.soft_delete_users() RETURNS TRIGGER
LANGUAGE plpgsql AS $$
DECLARE
    current_user_id UUID;
BEGIN
    current_user_id := COALESCE(
        NULLIF(current_setting('app.current_user_id', true), '')::UUID,
        erp.system_user_id()
    );

    UPDATE erp.users_table
    SET username = OLD.username || ':' || OLD.id,
        email = OLD.email || ':' || OLD.id,
        deleted_at = CURRENT_TIMESTAMP,
        deleted_by_id = current_user_id
    WHERE id = OLD.id;

    RETURN NULL;
END; $$;


CREATE OR REPLACE TRIGGER soft_delete_users
INSTEAD OF DELETE
ON erp.users
FOR EACH ROW
EXECUTE PROCEDURE erp.soft_delete_users();
