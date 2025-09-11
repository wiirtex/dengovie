package dengovie

import (
	mockDebtsService "dengovie/internal/mocks/service/debts_mock"
	mockTelegram "dengovie/internal/mocks/service/telegram_mock"
	mockUsersService "dengovie/internal/mocks/service/users_mock"
	mockStore "dengovie/internal/mocks/store"
	environment "dengovie/internal/utils/env"
	"testing"
)

func TestMain(m *testing.M) {
	environment.InitEnvs(map[environment.Key]string{
		environment.KeyJwtToken:           "test_jwt",
		environment.KeyPostgresConnString: "test_conn_string",
		environment.KeyTelegramBotToken:   "test_token",
	})

	m.Run()
}

type env struct {
	storage      *mockStore.MockStorage
	debtsService *mockDebtsService.MockService
	usersService *mockUsersService.MockService
	sender       *mockTelegram.MockClient
}

func newEnv(t *testing.T) *env {
	return &env{
		storage:      mockStore.NewMockStorage(t),
		debtsService: mockDebtsService.NewMockService(t),
		usersService: mockUsersService.NewMockService(t),
		sender:       mockTelegram.NewMockClient(t),
	}
}

func newController(e *env) *Controller {
	return NewController(
		e.storage,
		e.debtsService,
		e.usersService,
		e.sender,
	)
}
