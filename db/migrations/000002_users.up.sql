CREATE TABLE IF NOT EXISTS erp.users_table (
    id UUID PRIMARY KEY DEFAULT uuidv4(),

    name TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL,
    email TEXT UNIQUE NOT NULL,

    first_name TEXT,
    last_name TEXT,

    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT current_timestamp,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT current_timestamp,
    deleted_at TIMESTAMP WITH TIME ZONE
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
VALUES ('system', '', '_'), ('admin', '', '');


CREATE OR REPLACE FUNCTION erp.prevent_default_users_modification() RETURNS TRIGGER
LANGUAGE plpgsql AS $$ BEGIN
    IF OLD.name IN ('admin', 'system') THEN
        IF NEW IS NULL THEN
            RAISE EXCEPTION 'admin and system users cannot be deleted';
        END IF;
        IF OLD.name <> NEW.name THEN
            RAISE EXCEPTION 'admin and system user names cannot be modified';
        END IF;
    END IF;
    RETURN NEW;
END; $$;


CREATE OR REPLACE TRIGGER prevent_default_users_modification
BEFORE UPDATE OR DELETE
ON erp.users_table
FOR EACH ROW
EXECUTE PROCEDURE erp.prevent_default_users_modification();


CREATE MATERIALIZED VIEW IF NOT EXISTS erp.users_system AS
SELECT * FROM erp.users_table
WHERE name = 'system'
WITH DATA;


CREATE OR REPLACE FUNCTION erp.system_user_id() RETURNS UUID
STABLE
LANGUAGE plpgsql AS $$ BEGIN
    return (SELECT id FROM erp.users_system LIMIT 1);
END; $$;


ALTER TABLE IF EXISTS erp.users_table
ADD COLUMN IF NOT EXISTS created_by_id UUID NOT NULL DEFAULT erp.system_user_id() REFERENCES erp.users_table (
    id
) ON UPDATE CASCADE ON DELETE RESTRICT,
ADD COLUMN IF NOT EXISTS updated_by_id UUID NOT NULL DEFAULT erp.system_user_id() REFERENCES erp.users_table (
    id
) ON UPDATE CASCADE ON DELETE RESTRICT,
ADD COLUMN IF NOT EXISTS deleted_by_id UUID REFERENCES erp.users_table (
    id
) ON UPDATE CASCADE ON DELETE RESTRICT;


CREATE INDEX IF NOT EXISTS active_users_name_idx ON erp.users_table (name)
WHERE deleted_at IS NULL AND deleted_by_id IS NULL;


CREATE INDEX IF NOT EXISTS active_users_email_idx ON erp.users_table (email)
WHERE deleted_at IS NULL AND deleted_by_id IS NULL;


CREATE OR REPLACE VIEW erp.users AS
SELECT * FROM erp.users_table
WHERE
    name <> 'system'
    AND deleted_at IS NULL
    AND deleted_by_id IS NULL
WITH CHECK OPTION;


CREATE OR REPLACE FUNCTION erp.update_created_updated_by_id_users() RETURNS TRIGGER
LANGUAGE plpgsql AS $$
DECLARE
    current_user_id UUID;
BEGIN
    current_user_id := current_setting('app.current_user_id', true)::UUID;

    IF current_user_id IS NULL THEN
        RAISE EXCEPTION 'app.current_user_id must be set before inserting or updating';
    END IF;

    IF OLD IS NULL THEN
        NEW.created_by_id = current_user_id;
    END IF;
    NEW.updated_by_id = current_user_id;

    RETURN NEW;
END; $$;


CREATE OR REPLACE TRIGGER update_created_updated_by_id_users
BEFORE INSERT OR UPDATE
ON erp.users
FOR EACH ROW
EXECUTE PROCEDURE erp.update_created_updated_by_id_users();


CREATE OR REPLACE FUNCTION erp.soft_delete_users() RETURNS TRIGGER
LANGUAGE plpgsql AS $$
DECLARE
    current_user_id UUID;
BEGIN
    current_user_id := current_setting('app.current_user_id', true)::UUID;

    IF current_user_id IS NULL THEN
        RAISE EXCEPTION 'app.current_user_id must be set before deleting';
    END IF;

    UPDATE erp.users_table
    SET deleted_at = CURRENT_TIMESTAMP,
        deleted_by_id = current_user_id
    WHERE OLD.id;

    RETURN NULL;
END; $$;


CREATE OR REPLACE TRIGGER soft_delete_users
INSTEAD OF DELETE
ON erp.users
FOR EACH ROW
EXECUTE PROCEDURE erp.soft_delete_users();
