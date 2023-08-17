
CREATE TABLE IF NOT EXISTS schema_migrations
(
    version bigint  not null
        primary key,
    dirty   boolean not null
);

CREATE TABLE IF NOT EXISTS users
(
    id         uuid                                not null
        primary key,
    login      varchar(50)                         not null
        unique,
    password   text                                not null,
    created_at timestamp default CURRENT_TIMESTAMP not null,
    deleted_at timestamp
);

CREATE TABLE IF NOT EXISTS  secrets
(
    id          uuid                                not null
        primary key,
    user_id     uuid
        references users
            on delete cascade,
    secret_data bytea,
    secret_name bytea,
    created_at  timestamp default CURRENT_TIMESTAMP not null,
    deleted_at  timestamp
);

