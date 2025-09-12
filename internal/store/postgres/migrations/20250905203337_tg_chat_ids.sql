-- +goose Up
-- +goose StatementBegin

alter table users
    add column chat_id bigint default 0;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

alter table users
    drop column chat_id;

-- +goose StatementEnd
