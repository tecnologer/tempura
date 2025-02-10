package telegram

import (
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	*tgbotapi.BotAPI
}

func NewBot(token string, verbose bool) (*Bot, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, fmt.Errorf("creating bot: %w", err)
	}

	bot.Debug = verbose

	log.Printf("authorized on account %s", bot.Self.UserName)

	return &Bot{bot}, nil
}

func (b *Bot) SendMessage(chatID int64, text string) error {
	msg := tgbotapi.NewMessage(chatID, text)

	_, err := b.Send(msg)
	if err != nil {
		return fmt.Errorf("sending message: %w", err)
	}

	return nil
}
