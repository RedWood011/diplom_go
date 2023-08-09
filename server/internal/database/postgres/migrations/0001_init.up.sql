
create table if not exists schema_migrations
(
    version bigint  not null,
    dirty   boolean not null,
    primary key (version)
    );

-- auto-generated definition
create table schema_migrations
(
    version bigint  not null
        primary key,
    dirty   boolean not null
);

alter table schema_migrations
    owner to qwerty;

create table users
(
    id         uuid                                not null
        primary key,
    login      varchar(50)                         not null
        unique,
    password   text                                not null,
    created_at timestamp default CURRENT_TIMESTAMP not null,
    deleted_at timestamp
);

alter table users
    owner to qwerty;

create table secrets
(
    id          uuid                                not null
        primary key,
    user_id     uuid
        references users
            on delete cascade,
    secret_data bytea,
    created_at  timestamp default CURRENT_TIMESTAMP not null,
    deleted_at  timestamp,
    secret_name bytea
);

alter table secrets
    owner to qwerty;
