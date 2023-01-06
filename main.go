package main

import (
	"log"
	"os"

	"github.com/TanLeYang/stickies/db"
	stickiesset "github.com/TanLeYang/stickies/stickies_set"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	dbConn, err := db.GetConnection()
	if err != nil {
		log.Fatalf("Failed to connect to db: %s", err)
	}
	db.PerformMigration(dbConn)

	stickiesSetRepo := stickiesset.StickiesSetDb{
		Db: dbConn,
	}

	tgBotToken := os.Getenv("TG_BOT_TOKEN")
	bot, err := tgbotapi.NewBotAPI(tgBotToken)
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = true

	log.Printf("Telegram Bot API Authorized on account %s", bot.Self.UserName)

	stickiesBot := NewStickiesBot(bot, &stickiesSetRepo)
	stickiesBot.Start()
}
