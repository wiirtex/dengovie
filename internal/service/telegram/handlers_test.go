package telegram

import (
	"bytes"
	"context"
	"database/sql"
	"dengovie/internal/domain"
	storageTypes "dengovie/internal/store/types"
	tg "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"log"
	"os"
	"testing"
)

func Test_client_messageHandler(t *testing.T) {
	type args struct {
		update *models.Update
	}
	tests := []struct {
		name         string
		args         args
		prepareMocks func(t *testing.T, e *env)
		wantLogs     []string // Ожидаемые логи (можно проверять через hook)
	}{
		{
			name: "OK - user found and chat ID matches",
			args: args{
				update: &models.Update{
					Message: &models.Message{
						Chat: models.Chat{
							ID:       100,
							Username: "testuser",
						},
					},
				},
			},
			prepareMocks: func(t *testing.T, e *env) {
				e.storage.On("GetUserByAlias", mock.Anything, storageTypes.GetUserByAliasInput{
					Alias: "testuser",
				}).Return(storageTypes.User{
					ID:     123,
					Name:   "Test User",
					Alias:  "testuser",
					ChatID: 100, // Chat ID совпадает
				}, nil)

				e.bot.On("SendMessage", mock.Anything, &tg.SendMessageParams{
					ChatID:    int64(100),
					Text:      "Я не обрабатываю запросы. Могу только писать сам",
					ParseMode: models.ParseModeMarkdownV1,
				}).Return(nil, nil)
			},
			wantLogs: []string{},
		},
		{
			name: "OK - user found but chat ID needs update",
			args: args{
				update: &models.Update{
					Message: &models.Message{
						Chat: models.Chat{
							ID:       200, // Новый chat ID
							Username: "testuser",
						},
					},
				},
			},
			prepareMocks: func(t *testing.T, e *env) {
				e.storage.On("GetUserByAlias", mock.Anything, storageTypes.GetUserByAliasInput{
					Alias: "testuser",
				}).Return(storageTypes.User{
					ID:     123,
					Name:   "Test User",
					Alias:  "testuser",
					ChatID: 100, // Старый chat ID
				}, nil)

				e.storage.On("UpdateUserChatID", mock.Anything, storageTypes.UpdateUserChatIDInput{
					UserID:    domain.UserID(123),
					NewChatID: int64(200),
				}).Return(nil)

				e.bot.On("SendMessage", mock.Anything, &tg.SendMessageParams{
					ChatID:    int64(200),
					Text:      "Я обновил свои списки. Теперь тебе будут приходить уведомления",
					ParseMode: models.ParseModeMarkdownV1,
				}).Return(nil, nil)
			},
			wantLogs: []string{"user chat ID mismatch: user 100, chat ID 200"},
		},
		{
			name: "FAIL - user not found in database",
			args: args{
				update: &models.Update{
					Message: &models.Message{
						Chat: models.Chat{
							ID:       100,
							Username: "nonexistent",
						},
					},
				},
			},
			prepareMocks: func(t *testing.T, e *env) {
				e.storage.On("GetUserByAlias", mock.Anything, storageTypes.GetUserByAliasInput{
					Alias: "nonexistent",
				}).Return(storageTypes.User{}, sql.ErrNoRows)

				e.bot.On("SendMessage", mock.Anything, &tg.SendMessageParams{
					ChatID:    int64(100),
					Text:      "Сначала нужно зарегистрироваться под своим алиасом nonexistent",
					ParseMode: models.ParseModeMarkdownV1,
				}).Return(nil, nil)
			},
			wantLogs: []string{"no user found: alias nonexistent"},
		},
		{
			name: "FAIL - storage error when getting user",
			args: args{
				update: &models.Update{
					Message: &models.Message{
						Chat: models.Chat{
							ID:       100,
							Username: "testuser",
						},
					},
				},
			},
			prepareMocks: func(t *testing.T, e *env) {
				e.storage.On("GetUserByAlias", mock.Anything, storageTypes.GetUserByAliasInput{
					Alias: "testuser",
				}).Return(storageTypes.User{}, assert.AnError)

				e.bot.On("SendMessage", mock.Anything, &tg.SendMessageParams{
					ChatID:    int64(100),
					Text:      "Что-то пошло не так. Попробуй позже",
					ParseMode: models.ParseModeMarkdownV1,
				}).Return(nil, assert.AnError)
			},
			wantLogs: []string{},
		},
		{
			name: "FAIL - error updating user chat ID",
			args: args{
				update: &models.Update{
					Message: &models.Message{
						Chat: models.Chat{
							ID:       200,
							Username: "testuser",
						},
					},
				},
			},
			prepareMocks: func(t *testing.T, e *env) {
				e.storage.On("GetUserByAlias", mock.Anything, storageTypes.GetUserByAliasInput{
					Alias: "testuser",
				}).Return(storageTypes.User{
					ID:     123,
					Name:   "Test User",
					Alias:  "testuser",
					ChatID: 100,
				}, nil)

				e.storage.On("UpdateUserChatID", mock.Anything, storageTypes.UpdateUserChatIDInput{
					UserID:    domain.UserID(123),
					NewChatID: int64(200),
				}).Return(assert.AnError)
			},
			wantLogs: []string{
				"user chat ID mismatch: user 100, chat ID 200",
				"failed to update user's chatID:",
			},
		},
		{
			name: "FAIL - error sending message after chat ID update",
			args: args{
				update: &models.Update{
					Message: &models.Message{
						Chat: models.Chat{
							ID:       200,
							Username: "testuser",
						},
					},
				},
			},
			prepareMocks: func(t *testing.T, e *env) {
				e.storage.On("GetUserByAlias", mock.Anything, storageTypes.GetUserByAliasInput{
					Alias: "testuser",
				}).Return(storageTypes.User{
					ID:     123,
					Name:   "Test User",
					Alias:  "testuser",
					ChatID: 100,
				}, nil)

				e.storage.On("UpdateUserChatID", mock.Anything, storageTypes.UpdateUserChatIDInput{
					UserID:    domain.UserID(123),
					NewChatID: int64(200),
				}).Return(nil)

				e.bot.On("SendMessage", mock.Anything, &tg.SendMessageParams{
					ChatID:    int64(200),
					Text:      "Я обновил свои списки. Теперь тебе будут приходить уведомления",
					ParseMode: models.ParseModeMarkdownV1,
				}).Return(nil, assert.AnError)
			},
			wantLogs: []string{
				"user chat ID mismatch: user 100, chat ID 200",
				"failed to send default message:",
			},
		},
		{
			name: "FAIL - error sending message for not found user",
			args: args{
				update: &models.Update{
					Message: &models.Message{
						Chat: models.Chat{
							ID:       100,
							Username: "nonexistent",
						},
					},
				},
			},
			prepareMocks: func(t *testing.T, e *env) {
				e.storage.On("GetUserByAlias", mock.Anything, storageTypes.GetUserByAliasInput{
					Alias: "nonexistent",
				}).Return(storageTypes.User{}, sql.ErrNoRows)

				e.bot.On("SendMessage", mock.Anything, &tg.SendMessageParams{
					ChatID:    int64(100),
					Text:      "Сначала нужно зарегистрироваться под своим алиасом nonexistent",
					ParseMode: models.ParseModeMarkdownV1,
				}).Return(nil, assert.AnError)
			},
			wantLogs: []string{
				"no user found: alias nonexistent",
				"no user found: messageHandler.sendMessage:",
			},
		},
		{
			name: "FAIL - error sending default message",
			args: args{
				update: &models.Update{
					Message: &models.Message{
						Chat: models.Chat{
							ID:       100,
							Username: "testuser",
						},
					},
				},
			},
			prepareMocks: func(t *testing.T, e *env) {
				e.storage.On("GetUserByAlias", mock.Anything, storageTypes.GetUserByAliasInput{
					Alias: "testuser",
				}).Return(storageTypes.User{
					ID:     123,
					Name:   "Test User",
					Alias:  "testuser",
					ChatID: 100,
				}, nil)

				e.bot.On("SendMessage", mock.Anything, &tg.SendMessageParams{
					ChatID:    int64(100),
					Text:      "Я не обрабатываю запросы. Могу только писать сам",
					ParseMode: models.ParseModeMarkdownV1,
				}).Return(nil, assert.AnError)
			},
			wantLogs: []string{"failed to send default message:"},
		},
		{
			name: "OK - user without username",
			args: args{
				update: &models.Update{
					Message: &models.Message{
						Chat: models.Chat{
							ID:       100,
							Username: "", // Пустой username
						},
					},
				},
			},
			prepareMocks: func(t *testing.T, e *env) {
				e.storage.On("GetUserByAlias", mock.Anything, storageTypes.GetUserByAliasInput{
					Alias: "",
				}).Return(storageTypes.User{}, sql.ErrNoRows)

				e.bot.On("SendMessage", mock.Anything, &tg.SendMessageParams{
					ChatID:    int64(100),
					Text:      "Сначала нужно зарегистрироваться под своим алиасом ",
					ParseMode: models.ParseModeMarkdownV1,
				}).Return(nil, nil)
			},
			wantLogs: []string{"no user found: alias "},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Создаем логгер с hook для перехвата логов
			logOutput := &bytes.Buffer{}
			log.SetOutput(logOutput)
			defer log.SetOutput(os.Stderr)

			e := newEnv(t)
			c := newClient(e)

			tt.prepareMocks(t, e)

			// Вызываем тестируемую функцию
			c.messageHandler(context.Background(), nil, tt.args.update)

			// Проверяем логи, если нужно
			if len(tt.wantLogs) > 0 {
				logs := logOutput.String()
				for _, expectedLog := range tt.wantLogs {
					assert.Contains(t, logs, expectedLog)
				}
			}

			// Проверяем моки
			e.storage.AssertExpectations(t)
			e.bot.AssertExpectations(t)
		})
	}
}

// Дополнительные тесты для edge cases
func Test_client_messageHandler_EdgeCases(t *testing.T) {

	t.Run("empty username with spaces", func(t *testing.T) {
		logOutput := &bytes.Buffer{}
		log.SetOutput(logOutput)
		defer log.SetOutput(os.Stderr)

		e := newEnv(t)
		c := newClient(e)

		e.storage.On("GetUserByAlias", mock.Anything, storageTypes.GetUserByAliasInput{
			Alias: "   ",
		}).Return(storageTypes.User{}, sql.ErrNoRows)

		e.bot.On("SendMessage", mock.Anything, &tg.SendMessageParams{
			ChatID:    int64(100),
			Text:      "Сначала нужно зарегистрироваться под своим алиасом    ",
			ParseMode: models.ParseModeMarkdownV1,
		}).Return(nil, nil)

		update := &models.Update{
			Message: &models.Message{
				Chat: models.Chat{
					ID:       100,
					Username: "   ", // Пробелы
				},
			},
		}

		c.messageHandler(context.Background(), nil, update)

		assert.Contains(t, logOutput.String(), "no user found: alias    ")
		e.storage.AssertExpectations(t)
		e.bot.AssertExpectations(t)
	})
}
