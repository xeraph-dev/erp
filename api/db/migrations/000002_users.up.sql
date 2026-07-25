CREATE TABLE IF NOT EXISTS auth.__users (
    id              UUID    PRIMARY KEY DEFAULT uuidv4(),

    username        TEXT    UNIQUE  NOT NULL,
    email           TEXT    UNIQUE  NOT NULL,
    password_hash   TEXT            NOT NULL,

    first_name      TEXT,
    last_name       TEXT,

    created_at      TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at      TIMESTAMP WITH TIME ZONE
);


INSERT INTO auth.__users (username, email, password_hash)
VALUES ('system', '', '');


CREATE OR REPLACE FUNCTION auth.system_user_id() RETURNS UUID
STABLE
LANGUAGE plpgsql AS $$ BEGIN
    RETURN (SELECT id FROM auth.__users WHERE username = 'system' LIMIT 1);
END; $$;


ALTER TABLE IF EXISTS auth.__users
ADD COLUMN IF NOT EXISTS created_by_id UUID NOT NULL DEFAULT auth.system_user_id()  REFERENCES auth.__users (id) ON UPDATE CASCADE ON DELETE RESTRICT,
ADD COLUMN IF NOT EXISTS updated_by_id UUID NOT NULL DEFAULT auth.system_user_id()  REFERENCES auth.__users (id) ON UPDATE CASCADE ON DELETE RESTRICT,
ADD COLUMN IF NOT EXISTS deleted_by_id UUID                                         REFERENCES auth.__users (id) ON UPDATE CASCADE ON DELETE RESTRICT;


CREATE UNIQUE INDEX IF NOT EXISTS active_users_username_idx
ON auth.__users (username)
WHERE deleted_at IS NULL AND deleted_by_id IS NULL;


CREATE UNIQUE INDEX IF NOT EXISTS active_users_email_idx
ON auth.__users (email)
WHERE deleted_at IS NULL AND deleted_by_id IS NULL;


CREATE OR REPLACE FUNCTION auth.current_user_id() RETURNS UUID
STABLE
LANGUAGE plpgsql AS $$ BEGIN
    RETURN COALESCE(
        NULLIF(current_setting('app.current_user_id', true), '')::UUID,
        auth.system_user_id()
    );
END; $$;


CREATE OR REPLACE FUNCTION auth.update_modification_fields() RETURNS TRIGGER
LANGUAGE plpgsql AS $$ BEGIN
    IF OLD IS NULL THEN
        NEW.created_at = CURRENT_TIMESTAMP;
        NEW.created_by_id := auth.current_user_id();
    END IF;
    NEW.updated_at = CURRENT_TIMESTAMP;
    NEW.updated_by_id := auth.current_user_id();

    RETURN NEW;
END; $$;


CREATE OR REPLACE TRIGGER update_modification_fields
BEFORE INSERT OR UPDATE
ON auth.__users
FOR EACH ROW
EXECUTE FUNCTION auth.update_modification_fields();


CREATE OR REPLACE VIEW auth.users AS
SELECT
    id,
    username,
    email,
    password_hash,
    first_name,
    last_name
FROM auth.__users
WHERE username <> 'system'
  AND deleted_at IS NULL
  AND deleted_by_id IS NULL
WITH CHECK OPTION;


CREATE OR REPLACE FUNCTION auth.soft_delete_users() RETURNS TRIGGER
LANGUAGE plpgsql AS $$ BEGIN
    UPDATE auth.__users
    SET username = OLD.username || ':' || OLD.id,
        email = OLD.email || ':' || OLD.id,
        deleted_at = CURRENT_TIMESTAMP,
        deleted_by_id = auth.current_user_id()
    WHERE id = OLD.id;

    RETURN NULL;
END; $$;


CREATE OR REPLACE TRIGGER soft_delete_users
INSTEAD OF DELETE
ON auth.users
FOR EACH ROW
EXECUTE FUNCTION auth.soft_delete_users();
