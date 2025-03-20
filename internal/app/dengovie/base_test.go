package dengovie

import (
	"dengovie/internal/mocks"
	"testing"
)

type env struct {
	storage      *mocks.MockStorage
	debtsService *mocks.MockDebtsService
}

func newEnv(t *testing.T) *env {
	return &env{
		storage:      mocks.NewMockStorage(t),
		debtsService: mocks.NewMockDebtsService(t),
	}
}

func newController(e *env) *Controller {
	return NewController(
		e.storage,
		e.debtsService,
	)
}
