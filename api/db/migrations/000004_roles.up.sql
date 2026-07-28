CREATE TABLE IF NOT EXISTS auth.__roles (
    id              UUID        PRIMARY KEY DEFAULT uuidv7(),
    role_name       TEXT        NOT NULL,
    role_level      SMALLINT    NOT NULL,

    created_at      TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at      TIMESTAMP WITH TIME ZONE,

    created_by_id   UUID        NOT NULL DEFAULT auth.system_user_id()  REFERENCES auth.__users (id) ON UPDATE CASCADE ON DELETE RESTRICT,
    updated_by_id   UUID        NOT NULL DEFAULT auth.system_user_id()  REFERENCES auth.__users (id) ON UPDATE CASCADE ON DELETE RESTRICT,
    deleted_by_id   UUID
);


INSERT INTO auth.__roles (role_name, role_level)
VALUES ('admin', 32767), ('user', 0);


CREATE OR REPLACE FUNCTION auth.admin_role_id() RETURNS UUID
STABLE
LANGUAGE plpgsql AS $$ BEGIN
    RETURN (SELECT id FROM auth.__roles WHERE role_name = 'admin' LIMIT 1);
END; $$;


CREATE OR REPLACE FUNCTION auth.user_role_id() RETURNS UUID
STABLE
LANGUAGE plpgsql AS $$ BEGIN
    RETURN (SELECT id FROM auth.__roles WHERE role_name = 'user' LIMIT 1);
END; $$;


CREATE UNIQUE INDEX IF NOT EXISTS active_roles_role_name_idx
ON auth.__roles (role_name)
WHERE deleted_at IS NULL AND deleted_by_id IS NULL;


CREATE OR REPLACE TRIGGER update_modification_fields
BEFORE INSERT OR UPDATE
ON auth.__roles
FOR EACH ROW
EXECUTE FUNCTION auth.update_modification_fields();


CREATE OR REPLACE VIEW auth.roles AS
SELECT id, role_name, role_level
FROM auth.__roles
WHERE deleted_at IS NULL
  AND deleted_by_id IS NULL
WITH CHECK OPTION;


CREATE OR REPLACE FUNCTION auth.soft_delete_roles() RETURNS TRIGGER
LANGUAGE plpgsql AS $$ BEGIN
    UPDATE auth.__roles
    SET role_name = OLD.role_name || ':' || OLD.id,
        deleted_at = CURRENT_TIMESTAMP,
        deleted_by_id = auth.current_user_id()
    WHERE id = OLD.id;

    return NULL;
END; $$;


CREATE OR REPLACE TRIGGER soft_delete_roles
INSTEAD OF DELETE
ON auth.roles
FOR EACH ROW
EXECUTE FUNCTION auth.soft_delete_roles();


CREATE TABLE IF NOT EXISTS auth.roles_users (
    role_id         UUID    NOT NULL    REFERENCES auth.__roles (id) ON UPDATE CASCADE ON DELETE CASCADE,
    user_id         UUID    NOT NULL    REFERENCES auth.__users (id) ON UPDATE CASCADE ON DELETE CASCADE,

    created_at      TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,

    created_by_id   UUID    NOT NULL DEFAULT auth.system_user_id()  REFERENCES auth.__users (id) ON UPDATE CASCADE ON DELETE RESTRICT,
    updated_by_id   UUID    NOT NULL DEFAULT auth.system_user_id()  REFERENCES auth.__users (id) ON UPDATE CASCADE ON DELETE RESTRICT,

    PRIMARY KEY (role_id, user_id)
);


CREATE OR REPLACE TRIGGER update_modification_fields
BEFORE INSERT OR UPDATE
ON auth.roles_users
FOR EACH ROW
EXECUTE PROCEDURE auth.update_modification_fields();


INSERT INTO auth.__users (username, email, password_hash)
VALUES ('admin', 'admin@erp.local', '');


INSERT INTO auth.roles_users (role_id, user_id)
VALUES (auth.admin_role_id(), (SELECT id FROM auth.users WHERE username = 'admin' LIMIT 1));
