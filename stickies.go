package main

import (
	"bufio"
	"context"
	"log"
	"os"

	"github.com/TanLeYang/stickies/config"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type StickiesBot struct {
	tgBotAPI *tgbotapi.BotAPI
	conf     config.StickiesConfig
}

func NewStickiesBot(tgBotAPI *tgbotapi.BotAPI, conf config.StickiesConfig) *StickiesBot {
	return &StickiesBot{
		tgBotAPI,
		conf,
	}
}

func (sb *StickiesBot) Start() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := sb.tgBotAPI.GetUpdatesChan(u)

	go sb.receiveUpdates(ctx, updates)

	log.Println("Start listening for updates, press enter to stop")

	// Wait for a newline symbol, then cancel handling updates
	bufio.NewReader(os.Stdin).ReadBytes('\n')
	cancel()
}

func (sb *StickiesBot) receiveUpdates(ctx context.Context, updates tgbotapi.UpdatesChannel) {
	for {
		select {
		case <-ctx.Done():
			return

		case update := <-updates:
			sb.handleUpdate(update)
		}
	}
}

func (sb *StickiesBot) handleUpdate(update tgbotapi.Update) {
	switch {
	// Handle messages
	case update.Message != nil:
		sb.handleMessage(update.Message)
		break

	// Handle button clicks
	case update.CallbackQuery != nil:
		log.Fatalf("Handling button clicks not implemented yet")
	}
}

func (sb *StickiesBot) handleMessage(message *tgbotapi.Message) {
	user := message.From
	text := message.Text

	if user == nil {
		return
	}

	log.Printf("%s: %s", user.FirstName, text)

	if text == "new-sticker" {
		sb.addStickerToPack()
	}

	reply := tgbotapi.NewMessage(message.Chat.ID, message.Text)
	reply.ReplyToMessageID = message.MessageID

	sb.tgBotAPI.Send(reply)
}

func (sb *StickiesBot) addStickerToPack() {
	// hard coded for now, in the future will take in a picture, upload it and get the id
	fileID := "CAACAgUAAxkBAAMUY67gccAgX8NFESBqsnPG-_aVNHoAAg0IAAJhGzhV3r2-56AaX0MtBA"

	// Hard coded sticker to upload to try
	addStickerConfig := tgbotapi.AddStickerConfig{
		UserID:       sb.conf.WalnutStickerUserID,
		Name:         sb.conf.WalnutStickerSetName,
		PNGSticker:   tgbotapi.FileID(fileID),
		TGSSticker:   nil,
		Emojis:       "🤩",
		MaskPosition: nil,
	}

	resp, err := sb.tgBotAPI.Request(addStickerConfig)
	if err != nil {
		log.Printf("Error making add sticker request: \n %s \n", err)
	} else {
		log.Printf("Response from make sticker request: \n %s \n", resp.Description)
	}
}

