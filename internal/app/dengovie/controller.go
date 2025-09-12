package dengovie

import (
	debtsTypes "dengovie/internal/service/debts/types"
	telegramTypes "dengovie/internal/service/telegram/types"
	usersTypes "dengovie/internal/service/users/types"
	storeTypes "dengovie/internal/store/types"
)

// Controller example
type Controller struct {
	storage      storeTypes.Storage
	debtsService debtsTypes.Service
	usersService usersTypes.Service
	sender       telegramTypes.Client
}

// NewController example
func NewController(
	storage storeTypes.Storage,
	debtsService debtsTypes.Service,
	usersService usersTypes.Service,
	sender telegramTypes.Client,
) *Controller {
	return &Controller{
		storage:      storage,
		debtsService: debtsService,
		usersService: usersService,
		sender:       sender,
	}
}
