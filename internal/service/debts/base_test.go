package debts

import (
	"dengovie/internal/mocks"
	"testing"
)

type env struct {
	mockStorage *mocks.MockStorage
}

func newEnv(t *testing.T) *env {
	return &env{
		mockStorage: mocks.NewMockStorage(t),
	}
}

func newService(e *env) *Service {
	return New(e.mockStorage)
}
