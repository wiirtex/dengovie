package telegram

import (
	"context"
	"database/sql"
	"dengovie/internal/service/telegram/types"
	storageTypes "dengovie/internal/store/types"
	"dengovie/internal/utils/env"
	"errors"
	"fmt"
	tg "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"log"
)

var _ types.Client = (*Client)(nil)

type Client struct {
	db  storageTypes.Storage
	bot *tg.Bot
}

func (c *Client) SendMessageToUserByAlias(ctx context.Context, alias string, message string) error {
	user, err := c.db.GetUserByAlias(ctx, storageTypes.GetUserByAliasInput{
		Alias: alias,
	})
	if err != nil {
		return fmt.Errorf("db.GetUserByAlias: %w", err)
	}

	err = c.sendMessage(ctx, user.ChatID, message)
	if err != nil {
		return fmt.Errorf("sendMessage: %w", err)
	}

	return nil
}

func NewClient(db storageTypes.Storage, tokens ...string) (*Client, error) {

	token := ""
	if len(tokens) != 0 {
		token = tokens[0]
	} else {
		var err error
		token, err = env.GetEnv(env.KeyTelegramBotToken)
		if err != nil {
			panic(fmt.Sprintf("can not get token: %v", err))
		}
	}

	fmt.Println(token)
	bot, err := tg.New(token)
	if err != nil {
		return nil, fmt.Errorf("can not init bot: %w", err)
	}

	c := &Client{
		db:  db,
		bot: bot,
	}

	bot.RegisterHandler(tg.HandlerTypeMessageText, "", tg.MatchTypePrefix, c.messageHandler)

	return c, nil
}

func (c *Client) Start(ctx context.Context) {
	log.Println("starting bot")
	go func() {
		c.bot.Start(ctx)
	}()
}

func (c *Client) messageHandler(ctx context.Context, _ *tg.Bot, update *models.Update) {

	fmt.Println("helllo")

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

func (c *Client) sendMessage(ctx context.Context, chatID int64, format string, args ...any) error {

	msg := format
	if len(args) != 0 {
		msg = fmt.Sprintf(format, args...)
	}

	fmt.Println(chatID)
	_, err := c.bot.SendMessage(ctx, &tg.SendMessageParams{
		ChatID:    chatID,
		Text:      msg,
		ParseMode: models.ParseModeMarkdownV1,
	})
	if err != nil {
		return fmt.Errorf("bot.SendMessage: %w", err)
	}

	return nil
}
