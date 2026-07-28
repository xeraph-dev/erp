CREATE TABLE IF NOT EXISTS auth.permissions (
    id              UUID    PRIMARY KEY DEFAULT uuidv7(),
    permission_name TEXT    UNIQUE NOT NULL,

    created_at      TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,

    created_by_id   UUID    NOT NULL DEFAULT auth.system_user_id()  REFERENCES auth.__users (id) ON UPDATE CASCADE ON DELETE RESTRICT,
    updated_by_id   UUID    NOT NULL DEFAULT auth.system_user_id()  REFERENCES auth.__users (id) ON UPDATE CASCADE ON DELETE RESTRICT,
);


CREATE OR REPLACE TRIGGER update_modification_fields
BEFORE INSERT OR UPDATE
ON auth.permissions
FOR EACH ROW
EXECUTE FUNCTION auth.update_modification_fields();


CREATE TABLE IF NOT EXISTS auth.permissions_roles (
    permission_id   UUID    NOT NULL REFERENCES auth.permissions (id) ON UPDATE CASCADE ON DELETE CASCADE,
    role_id         UUID    NOT NULL REFERENCES auth.__roles (id) ON UPDATE CASCADE ON DELETE CASCADE,

    created_at      TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,

    created_by_id   UUID    NOT NULL DEFAULT auth.system_user_id()  REFERENCES auth.__users (id) ON UPDATE CASCADE ON DELETE RESTRICT,
    updated_by_id   UUID    NOT NULL DEFAULT auth.system_user_id()  REFERENCES auth.__users (id) ON UPDATE CASCADE ON DELETE RESTRICT,

    PRIMARY KEY (permission_id, role_id)
);


CREATE OR REPLACE TRIGGER update_modification_fields
BEFORE INSERT OR UPDATE
ON auth.permissions_roles
FOR EACH ROW
EXECUTE FUNCTION auth.update_modification_fields();


CREATE TABLE IF NOT EXISTS auth.permissions_users (
    permission_id   UUID    NOT NULL REFERENCES auth.permissions (id) ON UPDATE CASCADE ON DELETE CASCADE,
    user_id         UUID    NOT NULL REFERENCES auth.__users (id) ON UPDATE CASCADE ON DELETE CASCADE,

    created_at      TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,

    created_by_id   UUID    NOT NULL DEFAULT auth.system_user_id()  REFERENCES auth.__users (id) ON UPDATE CASCADE ON DELETE RESTRICT,
    updated_by_id   UUID    NOT NULL DEFAULT auth.system_user_id()  REFERENCES auth.__users (id) ON UPDATE CASCADE ON DELETE RESTRICT,

    PRIMARY KEY (permission_id, user_id)
);


CREATE OR REPLACE TRIGGER update_modification_fields
BEFORE INSERT OR UPDATE
ON auth.permissions_users
FOR EACH ROW
EXECUTE FUNCTION auth.update_modification_fields();
