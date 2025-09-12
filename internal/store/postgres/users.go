package postgres

import (
	"context"
	"database/sql"
	"dengovie/internal/store/types"
	"fmt"
)

func (r *Repo) GetUserByAlias(ctx context.Context, input types.GetUserByAliasInput) (types.User, error) {

	query := `
select u.id, u.name, u.alias, u.chat_id
from users u
where u.alias = $1
`

	row := r.db.QueryRowContext(ctx, query, input.Alias)
	if row.Err() != nil {
		return types.User{}, fmt.Errorf("db.QueryContext: %w", row.Err())
	}

	var user types.User
	err := row.Scan(&user.ID, &user.Name, &user.Alias, &user.ChatID)
	if err != nil {
		return user, fmt.Errorf("row.Scan: %w", err)
	}

	return user, nil
}

func (r *Repo) GetUserByID(ctx context.Context, input types.GetUserByIDInput) (types.User, error) {

	query := `
select u.id, u.name, u.alias, u.chat_id
from users u
where u.id = $1
`

	row := r.db.QueryRowContext(ctx, query, input.UserID)
	if row.Err() != nil {
		return types.User{}, fmt.Errorf("db.QueryContext: %w", row.Err())
	}

	var user types.User
	err := row.Scan(&user.ID, &user.Name, &user.Alias, &user.ChatID)
	if err != nil {
		return user, fmt.Errorf("row.Scan: %w", err)
	}

	return user, nil
}

func (r *Repo) UpdateUserName(ctx context.Context, input types.UpdateUserNameInput) error {

	query := `
update users
set name = $2
where id = $1
`

	_, err := r.db.ExecContext(ctx, query, input.UserID, input.NewName)
	if err != nil {
		return fmt.Errorf("db.ExecContext: %w", err)
	}

	return nil
}

// DeleteUser монстр запрос, который удаляет юзера из всех таблиц. Начисто. Без проверок.
func (r *Repo) DeleteUser(ctx context.Context, input types.DeleteUserInput) (errTx error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		errTx = fmt.Errorf("db.BeginTx: %w", err)
		return errTx
	}

	defer func(tx *sql.Tx) {
		if errTx != nil {
			err := tx.Rollback()
			if err != nil {
				errTx = fmt.Errorf("can not rollback trasaction: %w, after getting error: %w", err, errTx)
				return
			}
		}

		err := tx.Commit()
		if err != nil {
			errTx = fmt.Errorf("can not commit trasaction: %w", err)
			return
		}
	}(tx)

	query := `
delete 
from users
where id = $1
`
	_, err = tx.ExecContext(ctx, query, input.UserID)
	if err != nil {
		errTx = fmt.Errorf("db.ExecContext on delete from users: %w", err)
		return errTx
	}

	query = `
delete 
from debts
where user_id = $1 or another_user_id = $1
`
	_, err = tx.ExecContext(ctx, query, input.UserID)
	if err != nil {
		errTx = fmt.Errorf("db.ExecContext on delete from debts: %w", err)
		return errTx
	}

	query = `
delete 
from user_groups
where user_id = $1
`
	_, err = tx.ExecContext(ctx, query, input.UserID)
	if err != nil {
		errTx = fmt.Errorf("db.ExecContext on delete from user_groups: %w", err)
		return errTx
	}

	return nil
}

func (r *Repo) UpdateUserChatID(ctx context.Context, input types.UpdateUserChatIDInput) error {

	query := `
update users
set chat_id = $2
where id = $1
`

	_, err := r.db.ExecContext(ctx, query, input.UserID, input.NewChatID)
	if err != nil {
		return fmt.Errorf("db.ExecContext: %w", err)
	}

	return nil
}

func (r *Repo) CreateUser(ctx context.Context, input types.CreateUserInput) (types.User, error) {
	query := `
insert into users (name, alias)
values ($1, $2)
returning id, name, alias, chat_id 
`
	row := r.db.QueryRowContext(ctx, query, input.Name, input.Alias)
	if row.Err() != nil {
		return types.User{}, fmt.Errorf("db.QueryContext: %w", row.Err())
	}

	var user types.User
	err := row.Scan(&user.ID, &user.Name, &user.Alias, &user.ChatID)
	if err != nil {
		return user, fmt.Errorf("row.Scan: %w", err)
	}

	return user, nil
}
