CREATE TABLE IF NOT EXISTS erp.roles_users (
    user_id UUID REFERENCES erp.users_table (id) ON UPDATE CASCADE ON DELETE CASCADE,
    role_id UUID REFERENCES erp.roles_table (id) ON UPDATE CASCADE ON DELETE CASCADE,

    PRIMARY KEY (user_id, role_id)
);
