-- +goose Up
-- +goose StatementBegin

create table groups
(
    id   bigint primary key,
    name text
);

create table users
(
    id    bigint primary key,
    name  text,
    alias text
);

create table user_groups
(
    id       bigint primary key,
    user_id  bigint,
    group_id bigint
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

drop table user_groups;

drop table users;

drop table groups;

-- +goose StatementEnd
