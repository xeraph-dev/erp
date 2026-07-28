INSERT INTO auth.roles_users (role_id, user_id)
VALUES (auth.user_role_id(), @user_id);
