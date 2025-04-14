package postgres

import (
	"context"
	"dengovie/internal/domain"
	"dengovie/internal/store/types"
	"fmt"

	"github.com/lib/pq"
)

func (r *Repo) IsUserInGroup(ctx context.Context, userID domain.UserID, groupID domain.GroupID) (bool, error) {

	query := `
select exists (select 1 from user_groups ug where ug.group_id = $1 and ug.user_id = $2)
`

	row := r.db.QueryRowContext(ctx, query, groupID, userID)
	if row.Err() != nil {
		return false, fmt.Errorf("db.QueryRowContext: %w", row.Err())
	}

	var exists bool
	if err := row.Scan(&exists); err != nil {
		return false, fmt.Errorf("row.Scan: %w", err)
	}

	return exists, nil
}

func (r *Repo) AreUsersInGroup(ctx context.Context, userIDs []domain.UserID, groupID domain.GroupID) (bool, error) {

	query := `
select $3 = (select count(distinct ug.id) from user_groups ug where ug.group_id = $1 and ug.user_id = any($2))
`

	row := r.db.QueryRowContext(ctx, query, groupID, pq.Array(userIDs), len(userIDs))
	if row.Err() != nil {
		return false, fmt.Errorf("db.QueryRowContext: %w", row.Err())
	}

	var exists bool
	if err := row.Scan(&exists); err != nil {
		return false, fmt.Errorf("row.Scan: %w", err)
	}

	return exists, nil
}

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
