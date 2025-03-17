package postgres

import (
	"context"
	"dengovie/internal/store/types"
	"fmt"
)

func (r *Repo) ListUserGroups(ctx context.Context, input types.ListUserGroupsInput) ([]types.Group, error) {

	query := `
select g.id, g.name
from groups g
inner join user_groups ug on g.id = ug.group_id
where user_id = $1
`

	rows, err := r.db.QueryContext(ctx, query, input.UserID)
	if err != nil {
		return nil, fmt.Errorf("db.QueryContext: %w", err)
	}

	groups := make([]types.Group, 0, 5)
	for rows.Next() {
		var group types.Group
		errScan := rows.Scan(&group.ID, &group.Name)
		if errScan != nil {
			return nil, fmt.Errorf("rows.Scan: %w", errScan)
		}
		fmt.Println("add", group)

		groups = append(groups, group)
	}

	return groups, nil
}

func (r *Repo) ListUsersInGroup(ctx context.Context, input types.ListUsersInGroupInput) ([]types.User, error) {

	query := `
select u.id, u.name
from users u
inner join user_groups ug on u.id = ug.user_id
where ug.group_id = $1
`

	fmt.Println(input.GroupID)
	rows, err := r.db.QueryContext(ctx, query, input.GroupID)
	if err != nil {
		return nil, fmt.Errorf("db.QueryContext: %w", err)
	}

	users := make([]types.User, 0, 5)
	for rows.Next() {
		var user types.User
		errScan := rows.Scan(&user.ID, &user.Name)
		if errScan != nil {
			return nil, fmt.Errorf("rows.Scan: %w", errScan)
		}
		fmt.Println("add", user)

		users = append(users, user)
	}

	return users, nil

}
