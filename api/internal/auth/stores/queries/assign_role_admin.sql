INSERT INTO auth.roles_users (role_id, user_id)
VALUES (auth.admin_role_id(), @user_id);
