create table if not exists organizations
(
    id         bigserial    not null
    constraint organizations_pkey
    primary key,
    created_at timestamp with time zone default now(),
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
                             name       varchar(250) not null,
    active     boolean,
    uuid       uuid         not null
    constraint organizations_uuid_key
    unique
    );

alter table organizations
    owner to vitae;

create index if not exists orgs_uuid
    on organizations (uuid);

create table if not exists roles
(
    id           bigserial    not null
    constraint roles_pkey
    primary key,
    access_level bigint,
    name         varchar(100) not null,
    active       boolean
    );

alter table roles
    owner to vitae;

create index if not exists roles_uuid
    on roles (name);

create table if not exists users
(
    id          bigserial    not null
    constraint users_pkey
    primary key,
    created_at  timestamp with time zone default now(),
    updated_at  timestamp with time zone,
    deleted_at  timestamp with time zone,
                              first_name  varchar(100) not null,
    last_name   varchar(100) not null,
    email       varchar(255) not null,
    username    text,
    mobile      varchar(50),
    phone       varchar(50),
    address     varchar(50),
    external_id text
    constraint users_external_id_key
    unique,
    uuid        uuid         not null
    constraint users_uuid_key
    unique,
    constraint users_external_id_uuid_key
    unique (external_id, uuid)
    );

alter table users
    owner to vitae;

create index if not exists users_uuid
    on users (uuid);

create index if not exists users_externalid
    on users (external_id);

create index if not exists users_email
    on users (email);

create table if not exists profiles
(
    id              bigserial not null
    constraint profiles_pkey
    primary key,
    created_at      timestamp with time zone default now(),
    updated_at      timestamp with time zone,
    deleted_at      timestamp with time zone,
                                  uuid            uuid      not null
                                  constraint profiles_uuid_key
                                  unique,
                                  user_id         bigint    not null
                                  constraint profiles_user_id_fkey
                                  references users,
                                  organization_id bigint    not null
                                  constraint profiles_organization_id_fkey
                                  references organizations,
                                  role_id         bigint    not null
                                  constraint profiles_role_id_fkey
                                  references roles,
                                  active          boolean   not null
                                  );

alter table profiles
    owner to vitae;

create index if not exists profiles_uuid
    on profiles (uuid);

create table if not exists invitations
(
    id              bigserial    not null
    constraint invitations_pkey
    primary key,
    created_at      timestamp with time zone default now(),
    updated_at      timestamp with time zone,
    deleted_at      timestamp with time zone,
    token_hash      text
    constraint invitations_token_hash_key
    unique,
    expires_at      timestamp with time zone,
                                  invitor_id      bigint
                                  constraint invitations_invitor_id_fkey
                                  references users,
                                  organization_id bigint
                                  constraint invitations_organization_id_fkey
                                  references organizations,
                                  email           varchar(255) not null,
    used            boolean                  default false
    );

alter table invitations
    owner to vitae;



alter function uuid_nil() owner to postgres;

create function uuid_ns_dns() returns uuid
    immutable
    strict
    parallel safe
    language c
as
$$
begin
-- missing source code
end;
$$;

alter function uuid_ns_dns() owner to postgres;

create function uuid_ns_url() returns uuid
    immutable
    strict
    parallel safe
    language c
as
$$
begin
-- missing source code
end;
$$;

alter function uuid_ns_url() owner to postgres;

create function uuid_ns_oid() returns uuid
    immutable
    strict
    parallel safe
    language c
as
$$
begin
-- missing source code
end;
$$;

alter function uuid_ns_oid() owner to postgres;

create function uuid_ns_x500() returns uuid
    immutable
    strict
    parallel safe
    language c
as
$$
begin
-- missing source code
end;
$$;

alter function uuid_ns_x500() owner to postgres;

create function uuid_generate_v1() returns uuid
    strict
    parallel safe
    language c
as
$$
begin
-- missing source code
end;
$$;

alter function uuid_generate_v1() owner to postgres;

create function uuid_generate_v1mc() returns uuid
    strict
    parallel safe
    language c
as
$$
begin
-- missing source code
end;
$$;

alter function uuid_generate_v1mc() owner to postgres;

create function uuid_generate_v3(namespace uuid, name text) returns uuid
    immutable
    strict
    parallel safe
    language c
as
$$
begin
-- missing source code
end;
$$;

alter function uuid_generate_v3(uuid, text) owner to postgres;

create function uuid_generate_v4() returns uuid
    strict
    parallel safe
    language c
as
$$
begin
-- missing source code
end;
$$;

alter function uuid_generate_v4() owner to postgres;

create function uuid_generate_v5(namespace uuid, name text) returns uuid
    immutable
    strict
    parallel safe
    language c
as
$$
begin
-- missing source code
end;
$$;

alter function uuid_generate_v5(uuid, text) owner to postgres;

