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

create table if not exists project_statuses
(
    id         bigserial not null
    constraint project_statuses_pkey
    primary key,
    created_at timestamp with time zone default now(),
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
                             uuid       uuid      not null
                             constraint project_statuses_uuid_key
                             unique,
                             name       text,
                             "order"    bigint
                             );

alter table project_statuses
    owner to vitae;

create index if not exists project_statuses_uuid
    on project_statuses (uuid);

create table if not exists project_types
(
    id         bigserial not null
    constraint project_types_pkey
    primary key,
    created_at timestamp with time zone default now(),
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
                             uuid       uuid      not null
                             constraint project_types_uuid_key
                             unique,
                             name       text,
                             "order"    bigint
                             );

alter table project_types
    owner to vitae;

create table if not exists project_complexities
(
    id         bigserial not null
    constraint project_complexities_pkey
    primary key,
    created_at timestamp with time zone default now(),
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
                             uuid       uuid      not null
                             constraint project_complexities_uuid_key
                             unique,
                             name       text,
                             weight     bigint
                             );

alter table project_complexities
    owner to vitae;

create table if not exists project_sizes
(
    id         bigserial not null
    constraint project_sizes_pkey
    primary key,
    created_at timestamp with time zone default now(),
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
                             uuid       uuid      not null
                             constraint project_sizes_uuid_key
                             unique,
                             name       text,
                             weight     bigint
                             );

alter table project_sizes
    owner to vitae;

create table if not exists projects
(
    id                  bigserial    not null
    constraint projects_pkey
    primary key,
    created_at          timestamp with time zone default now(),
    updated_at          timestamp with time zone,
    deleted_at          timestamp with time zone,
                                      uuid                uuid         not null
                                      constraint projects_uuid_key
                                      unique,
                                      name                varchar(150) not null,
    description         text,
    rgt                 text,
    status_id           bigint       not null
    constraint projects_status_id_fkey
    references project_statuses,
    type_id             bigint       not null
    constraint projects_type_id_fkey
    references project_types,
    complexity_id       bigint       not null
    constraint projects_complexity_id_fkey
    references project_complexities,
    size_id             bigint       not null
    constraint projects_size_id_fkey
    references project_sizes,
    organization_id     bigint
    constraint projects_organization_id_fkey
    references organizations,
    est_start_date      date,
    est_end_date        date,
    act_start_date      date,
    act_end_date        date,
    open_for_time_entry boolean
    );

alter table projects
    owner to vitae;

create table if not exists strategic_alignments
(
    id              bigserial    not null
    constraint strategic_alignments_pkey
    primary key,
    created_at      timestamp with time zone default now(),
    updated_at      timestamp with time zone,
    deleted_at      timestamp with time zone,
                                  uuid            uuid         not null
                                  constraint strategic_alignments_uuid_key
                                  unique,
                                  name            varchar(100) not null,
    display_order   bigint,
    organization_id bigint
    constraint strategic_alignments_organization_id_fkey
    references organizations
    );

alter table strategic_alignments
    owner to vitae;

create table if not exists project_strategic_alignments
(
    project_id             bigserial not null
    constraint project_strategic_alignments_project_id_fkey
    references projects,
    strategic_alignment_id bigserial not null
    constraint project_strategic_alignments_strategic_alignment_id_fkey
    references strategic_alignments,
    constraint project_strategic_alignments_pkey
    primary key (project_id, strategic_alignment_id)
    );

alter table project_strategic_alignments
    owner to vitae;

create function uuid_nil() returns uuid
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

