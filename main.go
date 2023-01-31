package main

import (
	"log"

	"github.com/TanLeYang/stickies/config"
	"github.com/TanLeYang/stickies/db"
	stickiesset "github.com/TanLeYang/stickies/stickies_set"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	appConf := config.LoadAppConfig()

	dbConn, err := db.GetConnection(appConf.DbConf)
	if err != nil {
		log.Fatalf("Failed to connect to db: %s", err)
	}
	db.PerformMigration(dbConn)

	stickiesSetRepo := stickiesset.StickiesSetDb{
		Db: dbConn,
	}

	bot, err := tgbotapi.NewBotAPI(appConf.TgBotToken)
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = true

	log.Printf("Telegram Bot API Authorized on account %s", bot.Self.UserName)

	stickiesBot := NewStickiesBot(bot, &stickiesSetRepo, appConf.BotName)
	stickiesBot.Start()
}
