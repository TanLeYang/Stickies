package command

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type Command interface {
	Start(*tgbotapi.Message)
	Handle(*tgbotapi.Message)
}
