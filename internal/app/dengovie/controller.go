package dengovie

import (
	debtsTypes "dengovie/internal/service/debts/types"
	storeTypes "dengovie/internal/store/types"
)

// Controller example
type Controller struct {
	storage      storeTypes.Storage
	debtsService debtsTypes.DebtsService
}

// NewController example
func NewController(
	storage storeTypes.Storage,
	debtsService debtsTypes.DebtsService,
) *Controller {
	return &Controller{
		storage:      storage,
		debtsService: debtsService,
	}
}

// Message example
type Message struct {
	Message string `json:"message" example:"message"`
}
