package main

import (
	"log"
	"os"

	"github.com/TanLeYang/stickies/config"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	tgBotToken := os.Getenv("TG_BOT_TOKEN")
	bot, err := tgbotapi.NewBotAPI(tgBotToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Telegram Bot API Authorized on account %s", bot.Self.UserName)

	stickiesBotConf, err := config.ReadStickiesConfig()
	if err != nil {
		log.Panic(err)
	}
	stickiesBot := NewStickiesBot(bot, stickiesBotConf)
	stickiesBot.Start()
}
