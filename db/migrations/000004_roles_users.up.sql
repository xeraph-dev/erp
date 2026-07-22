CREATE TABLE IF NOT EXISTS erp.roles_users (
    role_id UUID
    REFERENCES erp.roles_table (id) ON UPDATE CASCADE ON DELETE CASCADE,
    user_id UUID
    REFERENCES erp.users_table (id) ON UPDATE CASCADE ON DELETE CASCADE,

    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,

    created_by_id UUID NOT NULL DEFAULT erp.system_user_id()
    REFERENCES erp.users_table (id) ON UPDATE CASCADE ON DELETE RESTRICT,
    updated_by_id UUID NOT NULL DEFAULT erp.system_user_id()
    REFERENCES erp.users_table (id) ON UPDATE CASCADE ON DELETE RESTRICT,

    PRIMARY KEY (role_id, user_id)
);


CREATE OR REPLACE TRIGGER update_modification_fields
BEFORE INSERT OR UPDATE
ON erp.roles_users
FOR EACH ROW
EXECUTE PROCEDURE erp.update_modification_fields();
