package dengovie

import (
	storeTypes "dengovie/internal/store/types"
)

// Controller example
type Controller struct {
	storage storeTypes.Storage
}

// NewController example
func NewController(storage storeTypes.Storage) *Controller {
	return &Controller{
		storage: storage,
	}
}

// Message example
type Message struct {
	Message string `json:"message" example:"message"`
}
