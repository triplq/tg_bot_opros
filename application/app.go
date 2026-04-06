package application

import (
	"parser/database"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type App struct {
	Model *database.Model
	Tg    *tgbotapi.BotAPI
}
