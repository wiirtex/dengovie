-- +goose Up
-- +goose StatementBegin


insert into users (name, alias)
values ('Тимур', 'wiirtex'),
       ('Вадим', 'galaxyshad'),
       ('Андрей', 'andrey'),
       ('Ярослав', 'yaroslav'),
       ('Диана', 'diana'),
       ('Васил', 'vasil'),
       ('Александр', 'alexander')
;

insert into groups (name)
values ('4119 group');

insert into user_groups (user_id, group_id)
select id, (select id from groups where name = '4119 group')
from users;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
