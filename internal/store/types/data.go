package types

import "dengovie/internal/domain"

type Group struct {
	ID   domain.GroupID
	Name string
}

type User struct {
	ID    domain.UserID
	Name  string
	Alias string
}

type UserDebt struct {
	AnotherUser User
	Amount      int64
}
