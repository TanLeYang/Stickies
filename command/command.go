package command

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type Command interface {
	// Handle messages and returns whether it expects additional messages
	Start(*tgbotapi.Message) bool
	Handle(*tgbotapi.Message) bool
}

type CommandStage string
