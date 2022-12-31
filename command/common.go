package command

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

func reply(botapi *tgbotapi.BotAPI, message *tgbotapi.Message, reply string) {
	r := tgbotapi.NewMessage(message.Chat.ID, reply)
	botapi.Send(r)
}
