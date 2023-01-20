package command

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type Command interface {
	// Handle messages and returns whether it expects additional messages
	Start(*tgbotapi.Message) CommandOngoingStatus
	Handle(*tgbotapi.Message) CommandOngoingStatus
}

type CommandStage string
type CommandOngoingStatus bool

const (
	CommandComplete CommandOngoingStatus = false
	CommandOngoing  CommandOngoingStatus = true
)
