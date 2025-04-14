package dengovie

import (
	debtsTypes "dengovie/internal/service/debts/types"
	usersTypes "dengovie/internal/service/users/types"
	storeTypes "dengovie/internal/store/types"
)

// Controller example
type Controller struct {
	storage      storeTypes.Storage
	debtsService debtsTypes.Service
	usersService usersTypes.Service
}

// NewController example
func NewController(
	storage storeTypes.Storage,
	debtsService debtsTypes.Service,
	usersService usersTypes.Service,
) *Controller {
	return &Controller{
		storage:      storage,
		debtsService: debtsService,
		usersService: usersService,
	}
}
