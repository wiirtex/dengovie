package dengovie

import (
	mockDebtsService "dengovie/internal/mocks/service/debts_mock"
	mockUsersService "dengovie/internal/mocks/service/users_mock"
	mockStore "dengovie/internal/mocks/store"
	environment "dengovie/internal/utils/env"
	"testing"
)

func TestMain(m *testing.M) {
	environment.InitEnvs(map[environment.Key]string{
		environment.KeyJwtToken: "test_jwt",
	})

	m.Run()
}

type env struct {
	storage      *mockStore.MockStorage
	debtsService *mockDebtsService.MockService
	usersService *mockUsersService.MockService
}

func newEnv(t *testing.T) *env {
	return &env{
		storage:      mockStore.NewMockStorage(t),
		debtsService: mockDebtsService.NewMockService(t),
		usersService: mockUsersService.NewMockService(t),
	}
}

func newController(e *env) *Controller {
	return NewController(
		e.storage,
		e.debtsService,
		e.usersService,
	)
}
