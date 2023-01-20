package interaction

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

func Reply(botapi *tgbotapi.BotAPI, message *tgbotapi.Message, reply string) {
	r := tgbotapi.NewMessage(message.Chat.ID, reply)
	botapi.Send(r)
}

func GenericErrorReply(botapi *tgbotapi.BotAPI, message *tgbotapi.Message) {
	Reply(botapi, message, "Sorry, something went wrong, please try again later!")
}
