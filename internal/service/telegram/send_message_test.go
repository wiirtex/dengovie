package telegram

import (
	"context"
	storageTypes "dengovie/internal/store/types"
	tg "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_client_sendMessage(t *testing.T) {
	type args struct {
		chatID int64
		format string
		args   []any
	}
	tests := []struct {
		name         string
		args         args
		prepareMocks func(t *testing.T, e *env)
		errWant      require.ErrorAssertionFunc
	}{
		{
			name: "OK",
			args: args{
				chatID: 100,
				format: "привет! Меня зовут %v",
				args:   []any{"Тимур"},
			},
			prepareMocks: func(t *testing.T, e *env) {
				e.bot.On("SendMessage", mock.Anything,
					&tg.SendMessageParams{
						ChatID:    int64(100),
						Text:      "привет! Меня зовут Тимур",
						ParseMode: models.ParseModeMarkdownV1,
					}).Return(nil, nil)
			},
			errWant: require.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			e := newEnv(t)
			c := newClient(e)

			tt.prepareMocks(t, e)

			err := c.sendMessage(context.Background(), tt.args.chatID, tt.args.format, tt.args.args...)
			tt.errWant(t, err)
		})
	}
}

func Test_client_SendMessageToUserByAlias(t *testing.T) {
	type args struct {
		alias   string
		message string
	}
	tests := []struct {
		name         string
		args         args
		prepareMocks func(t *testing.T, e *env)
		errWant      require.ErrorAssertionFunc
	}{
		{
			name: "OK",
			args: args{
				alias:   "ali",
				message: "привет",
			},
			prepareMocks: func(t *testing.T, e *env) {
				e.storage.On("GetUserByAlias", mock.Anything, storageTypes.GetUserByAliasInput{
					Alias: "ali",
				}).Return(storageTypes.User{
					ID:     101231,
					Name:   "Али",
					Alias:  "ali",
					ChatID: 100,
				}, nil)

				e.bot.On("SendMessage", mock.Anything,
					&tg.SendMessageParams{
						ChatID:    int64(100),
						Text:      "привет",
						ParseMode: models.ParseModeMarkdownV1,
					}).Return(nil, nil)
			},
			errWant: require.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			e := newEnv(t)
			c := newClient(e)

			tt.prepareMocks(t, e)

			err := c.SendMessageToUserByAlias(context.Background(), tt.args.alias, tt.args.message)
			tt.errWant(t, err)
		})
	}
}
