-- +goose Up
-- +goose StatementBegin

create table groups
(
    id   bigserial primary key,
    name text
);

create table users
(
    id    bigserial primary key,
    name  text,
    alias text
);

create table user_groups
(
    id       bigserial primary key,
    user_id  bigint,
    group_id bigint
);

create index user_groups_user_id_idx on user_groups(user_id);
create index user_groups_group_id_idx on user_groups(group_id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

drop table user_groups;

drop table users;

drop table groups;

-- +goose StatementEnd
