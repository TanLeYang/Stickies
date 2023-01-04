package main

import (
	"log"
	"os"

	"github.com/TanLeYang/stickies/config"
	"github.com/TanLeYang/stickies/db"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	_, err := db.GetConnection()
	if err != nil {
		log.Fatalf("Failed to connect to db: %s", err)
	}

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
