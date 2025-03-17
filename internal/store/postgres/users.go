package postgres

import (
	"context"
	"dengovie/internal/store/types"
	"fmt"
)

func (r *Repo) GetUserIDByAlias(ctx context.Context, input types.GetUserIDByAliasInput) (types.User, error) {

	query := `
select u.id, u.name
from users u
where u.alias = $1
`

	row := r.db.QueryRowContext(ctx, query, input.Alias)
	if row.Err() != nil {
		return types.User{}, fmt.Errorf("db.QueryContext: %w", row.Err())
	}

	var user types.User
	err := row.Scan(&user.ID, &user.Name)
	if err != nil {
		return user, fmt.Errorf("row.Scan: %w", err)
	}

	return user, nil
}
