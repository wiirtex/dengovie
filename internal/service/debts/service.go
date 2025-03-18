package debts

import "dengovie/internal/store/types"

type Service struct {
	storage types.Storage
}

func New(storage types.Storage) *Service {
	return &Service{
		storage: storage,
	}
}
