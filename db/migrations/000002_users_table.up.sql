CREATE TABLE IF NOT EXISTS erp.users_table (
    id UUID PRIMARY KEY DEFAULT uuidv4(),

    username TEXT NOT NULL,
    password_hash TEXT NOT NULL,
    email TEXT NOT NULL,

    first_name TEXT,
    last_name TEXT,

    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT current_timestamp,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT current_timestamp,
    deleted_at TIMESTAMP WITH TIME ZONE
);


INSERT INTO erp.users_table (username, password_hash, email)
VALUES ('system', '', '');


CREATE OR REPLACE FUNCTION erp.system_user_id() RETURNS UUID
STABLE
LANGUAGE plpgsql AS $$ BEGIN
    return (SELECT id FROM erp.users_table WHERE username = 'system' LIMIT 1);
END; $$;


ALTER TABLE IF EXISTS erp.users_table
ADD COLUMN IF NOT EXISTS created_by_id UUID NOT NULL
DEFAULT erp.system_user_id()
REFERENCES erp.users_table (id) ON UPDATE CASCADE ON DELETE RESTRICT,
ADD COLUMN IF NOT EXISTS updated_by_id UUID NOT NULL
DEFAULT erp.system_user_id()
REFERENCES erp.users_table (id) ON UPDATE CASCADE ON DELETE RESTRICT,
ADD COLUMN IF NOT EXISTS deleted_by_id UUID
REFERENCES erp.users_table (id) ON UPDATE CASCADE ON DELETE RESTRICT;


CREATE UNIQUE INDEX IF NOT EXISTS active_users_username_idx
ON erp.users_table (username)
WHERE deleted_at IS NULL AND deleted_by_id IS NULL;


CREATE UNIQUE INDEX IF NOT EXISTS active_users_email_idx
ON erp.users_table (email)
WHERE deleted_at IS NULL AND deleted_by_id IS NULL;


CREATE OR REPLACE FUNCTION erp.update_modification_fields()
RETURNS TRIGGER
LANGUAGE plpgsql AS $$
DECLARE
    current_user_id UUID;
BEGIN
    current_user_id := COALESCE(
        NULLIF(current_setting('app.current_user_id', true), '')::UUID,
        erp.system_user_id()
    );

    IF OLD IS NULL THEN
        NEW.created_at = CURRENT_TIMESTAMP;
        NEW.created_by_id = current_user_id;
    END IF;
    NEW.updated_at = CURRENT_TIMESTAMP;
    NEW.updated_by_id = current_user_id;

    RETURN NEW;
END; $$;


CREATE OR REPLACE TRIGGER update_modification_fields
BEFORE INSERT OR UPDATE
ON erp.users_table
FOR EACH ROW
EXECUTE PROCEDURE erp.update_modification_fields();
