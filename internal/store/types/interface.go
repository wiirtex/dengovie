package types

import (
	"context"
	"dengovie/internal/domain"
)

type Storage interface {
	IsUserInGroup(ctx context.Context, userID domain.UserID, groupID domain.GroupID) (bool, error)
	AreUsersInGroup(ctx context.Context, userIDs []domain.UserID, groupID domain.GroupID) (bool, error)
	ListUserGroups(ctx context.Context, input ListUserGroupsInput) ([]Group, error)
	ListUsersInGroup(ctx context.Context, input ListUsersInGroupInput) ([]User, error)

	GetUserIDByAlias(ctx context.Context, input GetUserIDByAliasInput) (User, error)

	ListUserDebts(ctx context.Context, input ListUserDebtsInput) ([]UserDebt, error)
	ShareDebt(ctx context.Context, input ShareDebtInput) error
	CreateEmptyDebts(ctx context.Context, input CreateEmptyDebtsInput) error

	PayDebt(ctx context.Context, input PayDebtInput) error
}
