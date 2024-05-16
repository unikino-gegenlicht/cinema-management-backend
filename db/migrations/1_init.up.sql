BEGIN;
CREATE SCHEMA "auth";

CREATE TABLE IF NOT EXISTS auth.users
(
    id                  uuid primary key default gen_random_uuid(),
    username            text not null,
    password            text not null,
    name                text not null,
    email               text not null,
    external            bool             default true,
    external_identifier text
);

CREATE TABLE IF NOT EXISTS auth.permissions
(
    id          uuid primary key default gen_random_uuid(),
    title       text not null,
    description text not null,
    scope       text not null,
    removable   bool             default false
);

CREATE TABLE IF NOT EXISTS auth.permission_groups
(
    id                   uuid primary key default gen_random_uuid(),
    title                text not null,
    description          text not null,
    included_permissions uuid[]
);

CREATE TABLE IF NOT EXISTS auth.permission_assignments
(
    "user"       uuid references auth.users (id) match full,
    permission uuid references auth.permissions (id),
    scope      int not null default 0
);
COMMIT;