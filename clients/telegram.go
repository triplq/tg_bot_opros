package clients

import (
	"fmt"
	"parser/config"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Init() (*tgbotapi.BotAPI, error) {
	bot, err := tgbotapi.NewBotAPI(config.Config("API_KEY"))
	if err != nil {
		return nil, fmt.Errorf("Error with tgbotapi")
	}

	bot.Debug = true
	return bot, nil
}
