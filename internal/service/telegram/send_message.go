package telegram

import (
	"context"
	storageTypes "dengovie/internal/store/types"
	"fmt"
	tg "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func (c *client) SendMessageToUserByAlias(ctx context.Context, alias string, message string) error {
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

func (c *client) sendMessage(ctx context.Context, chatID int64, format string, args ...any) error {

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
