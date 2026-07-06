CREATE TABLE IF NOT EXISTS erp.users_table (
    LIKE erp.base_table INCLUDING ALL,

    name        TEXT    UNIQUE NOT NULL,
    password    TEXT    NOT NULL,
    email       TEXT    UNIQUE NOT NULL,

    first_name  TEXT,
    last_name   TEXT
);


CREATE OR REPLACE FUNCTION erp.update_updated_at() RETURNS TRIGGER
    LANGUAGE plpgsql AS $$ BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END; $$;


CREATE OR REPLACE TRIGGER update_updated_at
    BEFORE UPDATE
    ON erp.users_table
    FOR EACH ROW
EXECUTE PROCEDURE erp.update_updated_at();


INSERT INTO erp.users_table (name, password, email)
VALUES ('system', '', '');


CREATE MATERIALIZED VIEW IF NOT EXISTS erp.users_system AS
SELECT * FROM erp.users_table
WHERE name = 'system'
WITH DATA;


CREATE OR REPLACE FUNCTION erp.system_user_id() RETURNS UUID
    IMMUTABLE
    LANGUAGE plpgsql AS $$ BEGIN
    return (SELECT id FROM erp.users_system LIMIT 1);
END; $$;


ALTER TABLE IF EXISTS erp.base_table
ADD COLUMN IF NOT EXISTS created_by_id UUID NOT NULL DEFAULT erp.system_user_id() REFERENCES erp.users_table (id) ON UPDATE CASCADE ON DELETE RESTRICT,
ADD COLUMN IF NOT EXISTS updated_by_id UUID NOT NULL DEFAULT erp.system_user_id() REFERENCES erp.users_table (id) ON UPDATE CASCADE ON DELETE RESTRICT,
ADD COLUMN IF NOT EXISTS deleted_by_id UUID REFERENCES erp.users_table (id) ON UPDATE CASCADE ON DELETE RESTRICT;


ALTER TABLE IF EXISTS erp.users_table
ADD COLUMN IF NOT EXISTS created_by_id UUID NOT NULL DEFAULT erp.system_user_id() REFERENCES erp.users_table (id) ON UPDATE CASCADE ON DELETE RESTRICT,
ADD COLUMN IF NOT EXISTS updated_by_id UUID NOT NULL DEFAULT erp.system_user_id() REFERENCES erp.users_table (id) ON UPDATE CASCADE ON DELETE RESTRICT,
ADD COLUMN IF NOT EXISTS deleted_by_id UUID REFERENCES erp.users_table (id) ON UPDATE CASCADE ON DELETE RESTRICT;


CREATE OR REPLACE FUNCTION erp.default_deleted_at() RETURNS TRIGGER
    LANGUAGE plpgsql AS $$ BEGIN
    IF NEW.deleted_at IS NULL AND NEW.deleted_by_id IS NOT NULL THEN
        NEW.deleted_at = CURRENT_TIMESTAMP;
    END IF;
    RETURN NEW;
END; $$;


CREATE OR REPLACE TRIGGER default_deleted_at
    BEFORE INSERT OR UPDATE
    ON erp.users_table
    FOR EACH ROW
EXECUTE PROCEDURE erp.default_deleted_at();


CREATE INDEX IF NOT EXISTS active_users_name_idx ON erp.users_table (name)
WHERE deleted_at IS NULL AND deleted_by_id IS NULL;


CREATE INDEX IF NOT EXISTS active_users_email_idx ON erp.users_table (email)
WHERE deleted_at IS NULL AND deleted_by_id IS NULL;


CREATE OR REPLACE VIEW erp.users AS
SELECT * FROM erp.users_table
WHERE deleted_at IS NULL
  AND deleted_by_id IS NULL;
