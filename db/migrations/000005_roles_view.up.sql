CREATE OR REPLACE VIEW erp.roles AS
SELECT
    id,
    role_name,
    role_level
FROM erp.roles_table
WHERE
    deleted_at IS NULL
    AND deleted_by_id IS NULL
WITH CHECK OPTION; -- noqa: disable=PRS


CREATE OR REPLACE FUNCTION erp.soft_delete_roles() RETURNS TRIGGER
LANGUAGE plpgsql AS $$
DECLARE
    current_user_id UUID;
BEGIN
    current_user_id := COALESCE(
        NULLIF(current_setting('app.current_user_id', true), '')::UUID,
        erp.system_user_id()
    );

    UPDATE erp.roles_table
    SET role_name = OLD.role_name || ':' || OLD.id,
        deleted_at = CURRENT_TIMESTAMP,
        deleted_by_id = current_user_id
    WHERE id = OLD.id;

    RETURN NULL;
END; $$;


CREATE OR REPLACE TRIGGER soft_delete_roles
INSTEAD OF DELETE
ON erp.roles
FOR EACH ROW
EXECUTE PROCEDURE erp.soft_delete_roles();
