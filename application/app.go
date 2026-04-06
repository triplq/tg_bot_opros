package application

import (
	"database/sql"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type App struct {
	DB *sql.DB
	tg *tgbotapi.BotAPI
}
