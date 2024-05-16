BEGIN;
DROP TABLE auth.permission_assignments CASCADE;
DROP TABLE auth.permission_groups CASCADE;
DROP TABLE auth.permissions CASCADE;
DROP TABLE auth.users CASCADE;
DROP SCHEMA auth;
COMMIT;