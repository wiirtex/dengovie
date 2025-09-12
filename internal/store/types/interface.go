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

	CreateUser(ctx context.Context, input CreateUserInput) (User, error)
	GetUserByAlias(ctx context.Context, input GetUserByAliasInput) (User, error)
	GetUserByID(ctx context.Context, input GetUserByIDInput) (User, error)
	UpdateUserName(ctx context.Context, input UpdateUserNameInput) error
	DeleteUser(ctx context.Context, input DeleteUserInput) error
	UpdateUserChatID(ctx context.Context, input UpdateUserChatIDInput) error

	ListUserDebts(ctx context.Context, input ListUserDebtsInput) ([]UserDebt, error)
	ShareDebt(ctx context.Context, input ShareDebtInput) error
	CreateEmptyDebts(ctx context.Context, input CreateEmptyDebtsInput) error
}
