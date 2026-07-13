CREATE TABLE IF NOT EXISTS erp.roles_table (
    id UUID PRIMARY KEY DEFAULT uuidv4(),

    role_name TEXT NOT NULL,
    role_level SMALLINT NOT NULL,

    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT current_timestamp,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT current_timestamp,
    deleted_at TIMESTAMP WITH TIME ZONE,

    created_by_id UUID NOT NULL DEFAULT erp.system_user_id()
    REFERENCES erp.users_table (id) ON UPDATE CASCADE ON DELETE RESTRICT,
    updated_by_id UUID NOT NULL DEFAULT erp.system_user_id()
    REFERENCES erp.users_table (id) ON UPDATE CASCADE ON DELETE RESTRICT,
    deleted_by_id UUID
    REFERENCES erp.users_table (id) ON UPDATE CASCADE ON DELETE RESTRICT
);


INSERT INTO erp.roles_table (role_name, role_level)
VALUES ('admin', 32767), ('user', 0);


CREATE UNIQUE INDEX IF NOT EXISTS active_roles_role_name_idx
ON erp.roles_table (role_name)
WHERE deleted_at IS NULL AND deleted_by_id IS NULL;


CREATE OR REPLACE TRIGGER update_modification_fields
BEFORE INSERT OR UPDATE
ON erp.roles_table
FOR EACH ROW
EXECUTE PROCEDURE erp.update_modification_fields();
