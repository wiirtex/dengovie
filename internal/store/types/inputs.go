package types

import "dengovie/internal/domain"

type ListUserGroupsInput struct {
	UserID domain.UserID
}

type ListUsersInGroupInput struct {
	GroupID domain.GroupID
}

type GetUserByAliasInput struct {
	Alias string
}

type GetUserByIDInput struct {
	UserID domain.UserID
}

type UpdateUserNameInput struct {
	UserID  domain.UserID
	NewName string
}

type DeleteUserInput struct {
	UserID domain.UserID
}

type ListUserDebtsInput struct {
	UserID domain.UserID
}

type ShareDebtInput struct {
	UserID           domain.UserID
	ChangeDebtAmount []ChangeUserDebtAmountInput
}

type ChangeUserDebtAmountInput struct {
	UserID domain.UserID
	Amount int64
}

type CreateEmptyDebtsInput struct {
	UserID         domain.UserID
	AnotherUserIDs []domain.UserID
}

type PayDebtInput struct {
	UserID         domain.UserID
	AnotherUserIDs []domain.UserID
	Amount         int64
}
