package telegram

import (
	"context"
	"dengovie/internal/service/telegram/types"
	storageTypes "dengovie/internal/store/types"
	environment "dengovie/internal/utils/env"
	"fmt"
	tg "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"log"
)

var _ types.Client = (*client)(nil)

type client struct {
	db  storageTypes.Storage
	bot bot
}

// обёртка над библиотечным типом, чтобы мокать
type bot interface {
	Start(ctx context.Context)
	SendMessage(ctx context.Context, params *tg.SendMessageParams) (*models.Message, error)
}

func NewClient(db storageTypes.Storage, tokens ...string) (*client, error) {

	token := ""
	if len(tokens) != 0 {
		token = tokens[0]
	} else {
		var err error
		token, err = environment.GetEnv(environment.KeyTelegramBotToken)
		if err != nil {
			panic(fmt.Sprintf("can not get token: %v", err))
		}
	}

	fmt.Println(token)
	bot, err := tg.New(token)
	if err != nil {
		return nil, fmt.Errorf("can not init bot: %w", err)
	}

	c := &client{
		db:  db,
		bot: bot,
	}

	bot.RegisterHandler(tg.HandlerTypeMessageText, "", tg.MatchTypePrefix, c.messageHandler)

	return c, nil
}

func (c *client) Start(ctx context.Context) {
	log.Println("starting bot")
	go func() {
		c.bot.Start(ctx)
	}()
}
