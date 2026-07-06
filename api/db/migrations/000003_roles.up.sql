CREATE TABLE IF NOT EXISTS erp.roles_table (
    LIKE erp.base_table INCLUDING ALL,

    name    TEXT        UNIQUE NOT NULL,
    level   SMALLINT    NOT NULL DEFAULT 0
);


CREATE OR REPLACE TRIGGER update_updated_at
    BEFORE INSERT OR UPDATE
    ON erp.roles_table
    FOR EACH ROW
EXECUTE PROCEDURE erp.update_updated_at();


CREATE OR REPLACE TRIGGER default_deleted_at
    BEFORE INSERT OR UPDATE
    ON erp.roles_table
    FOR EACH ROW
EXECUTE PROCEDURE erp.default_deleted_at();


CREATE INDEX IF NOT EXISTS active_users_roles_idx ON erp.roles_table (name)
WHERE deleted_at IS NULL AND deleted_by_id IS NULL;


CREATE OR REPLACE VIEW erp.roles AS
SELECT * FROM erp.roles_table
WHERE deleted_at IS NULL
  AND deleted_by_id IS NULL;


INSERT INTO erp.roles_table (name, level)
VALUES ('admin', 32767), ('user', 0);


CREATE OR REPLACE FUNCTION erp.prevent_default_roles_modification() RETURNS TRIGGER
    LANGUAGE plpgsql AS $$ BEGIN
    IF OLD.name IN ('admin', 'user') THEN
        RAISE EXCEPTION 'admin and user roles cannot be modified or deleted';
    END IF;
    RETURN NEW;
END; $$;


CREATE OR REPLACE TRIGGER prevent_default_roles_modification
    BEFORE UPDATE OR DELETE
    ON erp.roles_table
    FOR EACH ROW
EXECUTE PROCEDURE erp.prevent_default_roles_modification();
