package types

type ListUserGroupsInput struct {
	UserID int64
}

type ListUsersInGroupInput struct {
	GroupID int64
}

type GetUserIDByAliasInput struct {
	Alias string
}
