# My notes

- the system has a default admin user
- the only way to create a new user is through the registration process
- admin role has all permissions
- it's no possible to modify anything related to the admin, only the password
- admin role is always the highest role and no one can assign it to another user
- all tables should have a `created_by`, `updated_by` and `deleted_by` if needed, expect the `users` table
- there is a `system` user that represents the system itself in the `*_by` fields
- each role can has multiple permissions
- each user can has multiple roles
- roles are just a way to name a group of permissions, the system only check the permission that the user have
- permissions are allow-only, if the permission is no present, it's assumed that does not have that permission
- the `created_at` field is cannot be changed
- by default the system includes three type of users:
  - system: to represent operations made by the system itself, there is not a system role
  - admin: the highest privileged user, with role `admin`
  - user: the new users of the system, with role `user`
- the default roles cannot be changed or deleted because they are used by the system
- permissions are readonly because they are used by the system, the user cannot modify nor delete a permission
