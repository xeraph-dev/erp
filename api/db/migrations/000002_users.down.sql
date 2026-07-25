DROP TRIGGER IF EXISTS soft_delete_users ON auth.users;
DROP FUNCTION IF EXISTS auth.soft_delete_users;


DROP VIEW IF EXISTS auth.users;


DROP TRIGGER IF EXISTS update_modification_fields ON auth.__users;
DROP FUNCTION IF EXISTS auth.update_modification_fields;


DROP FUNCTION IF EXISTS auth.current_user_id;


DROP INDEX IF EXISTS active_users_email_idx;
DROP INDEX IF EXISTS active_users_username_idx;


ALTER TABLE IF EXISTS auth.__users
DROP COLUMN IF EXISTS created_by_id,
DROP COLUMN IF EXISTS updated_by_id,
DROP COLUMN IF EXISTS deleted_by_id;


DROP FUNCTION IF EXISTS auth.system_user_id;


DELETE FROM auth.__users
WHERE username = 'system';


DROP TABLE IF EXISTS auth.__users;
