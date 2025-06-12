-- +goose Up
-- +goose StatementBegin

create table debts
(
    id              bigserial primary key,
    user_id         bigint,
    another_user_id bigint,
    direction       int2,
    amount          bigint
);

create unique index debts_buyer_id_idx on debts (user_id, another_user_id, direction);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

drop table debts;

-- +goose StatementEnd
