package types

import "context"

type Storage interface {
	ListUserGroups(ctx context.Context, input ListUserGroupsInput) ([]Group, error)
	ListUsersInGroup(ctx context.Context, input ListUsersInGroupInput) ([]User, error)

	GetUserIDByAlias(ctx context.Context, input GetUserIDByAliasInput) (User, error)
}
