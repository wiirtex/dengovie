package telegram

import (
	"context"
	"database/sql"
	storageTypes "dengovie/internal/store/types"
	"errors"
	tg "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"log"
)

func (c *client) messageHandler(ctx context.Context, _ *tg.Bot, update *models.Update) {

	chatID := update.Message.Chat.ID
	userName := update.Message.Chat.Username

	user, err := c.db.GetUserByAlias(ctx, storageTypes.GetUserByAliasInput{
		Alias: userName,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("no user found: alias %v", userName)

			errSendMessage := c.sendMessage(ctx, update.Message.Chat.ID,
				"Сначала нужно зарегистрироваться под своим алиасом %v",
				update.Message.Chat.Username)
			if errSendMessage != nil {
				log.Printf("no user found: messageHandler.sendMessage: %v", errSendMessage)
			}
			return
		}

		log.Printf("failed to get user: %v", err)
		err = c.sendMessage(ctx,
			update.Message.Chat.ID,
			"Что-то пошло не так. Попробуй позже")
		if err != nil {
			log.Printf("failed to send message to user: %v", err)
		}
		return
	}

	if user.ChatID != chatID {
		log.Printf("user chat ID mismatch: user %v, chat ID %v", user.ChatID, chatID)
		errUpdateUser := c.db.UpdateUserChatID(ctx, storageTypes.UpdateUserChatIDInput{
			UserID:    user.ID,
			NewChatID: chatID,
		})
		if errUpdateUser != nil {
			log.Printf("failed to update user's chatID: %v", errUpdateUser)
			return
		}

		err = c.sendMessage(ctx,
			update.Message.Chat.ID,
			"Я обновил свои списки. Теперь тебе будут приходить уведомления")
		if err != nil {
			log.Printf("failed to send default message: %v", err)
		}
		return
	}

	err = c.sendMessage(ctx,
		update.Message.Chat.ID,
		"Я не обрабатываю запросы. Могу только писать сам")
	if err != nil {
		log.Printf("failed to send default message: %v", err)
	}
}
