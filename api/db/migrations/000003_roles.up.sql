CREATE TABLE IF NOT EXISTS erp.roles_table (
    LIKE erp.base_table INCLUDING ALL,

    name    TEXT    UNIQUE NOT NULL
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


CREATE INDEX IF NOT EXISTS active_roles_name_idx ON erp.roles_table (name)
WHERE deleted_at IS NULL AND deleted_by_id IS NULL;


CREATE OR REPLACE VIEW erp.roles AS
SELECT * FROM erp.roles_table
WHERE deleted_at IS NULL
  AND deleted_by_id IS NULL;
