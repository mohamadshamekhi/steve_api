package telegrambot

import (
	"context"
	"log"

	"github.com/go-telegram/bot"
)

type Config struct {
	AdminApiTokenBot     string `koanf:"admin_api_token_bot"`
	AdminChatID          int64  `koanf:"admin_chat_id"`
	DeveloperApiTokenBot string `koanf:"developer_api_token_bot"`
	DeveloperChatID      int64  `koanf:"developer_chat_id"`
}

type Adapter struct {
	cfg Config
	bot *bot.Bot
}

func (a *Adapter) GetBot() *bot.Bot {
	return a.bot
}

func NewAdmin(cfg Config) *Adapter {

	b, err := bot.New(cfg.AdminApiTokenBot)
	if err != nil {
		log.Printf("Failed to create telegram bot config: %v", err)
	}
	return &Adapter{
		cfg: cfg,
		bot: b,
	}
}

func (a *Adapter) SendMessageAdmin(ctx context.Context, text string) error {
	_, err := a.bot.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: a.cfg.AdminChatID,
		Text:   text,
	})
	return err
}

func NewDeveloper(cfg Config) *Adapter {
	b, err := bot.New(cfg.DeveloperApiTokenBot)
	if err != nil {
		log.Printf("Failed to create telegram bot config: %v", err)
	}
	return &Adapter{
		cfg: cfg,
		bot: b,
	}
}

func (a *Adapter) SendMessageDeveloper(ctx context.Context, text string) error {
	_, err := a.bot.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: a.cfg.DeveloperChatID,
		Text:   text,
	})
	return err
}
