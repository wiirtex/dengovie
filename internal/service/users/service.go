package users

import (
	"dengovie/internal/service/users/types"
	storeTypes "dengovie/internal/store/types"
)

type Service struct {
	storage storeTypes.Storage
}

func New(storage storeTypes.Storage) types.Service {
	return &Service{
		storage: storage,
	}
}
