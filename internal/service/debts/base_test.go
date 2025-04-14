package debts

import (
	mockStore "dengovie/internal/mocks/store"
	"testing"
)

type env struct {
	mockStorage *mockStore.MockStorage
}

func newEnv(t *testing.T) *env {
	return &env{
		mockStorage: mockStore.NewMockStorage(t),
	}
}

func newService(e *env) *Service {
	return New(e.mockStorage)
}
