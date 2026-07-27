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
