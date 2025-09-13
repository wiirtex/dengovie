package telegram

import (
	mockBot "dengovie/internal/mocks/service/telegram_internal_mock"
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
	storage *mockStore.MockStorage
	bot     *mockBot.Mockbot
}

func newEnv(t *testing.T) *env {
	return &env{
		storage: mockStore.NewMockStorage(t),
		bot:     mockBot.NewMockbot(t),
	}
}

func newClient(e *env) *client {
	return &client{
		db:  e.storage,
		bot: e.bot,
	}
}
