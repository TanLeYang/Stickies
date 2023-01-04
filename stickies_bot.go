package main

import (
	"bufio"
	"context"
	"log"
	"os"

	"github.com/TanLeYang/stickies/command"
	"github.com/TanLeYang/stickies/config"
	"github.com/TanLeYang/stickies/interaction"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type StickiesBot struct {
	tgBotAPI        *tgbotapi.BotAPI
	conf            config.StickiesConfig
	currrentCommand command.Command
}

func NewStickiesBot(tgBotAPI *tgbotapi.BotAPI, conf config.StickiesConfig) *StickiesBot {
	return &StickiesBot{
		tgBotAPI,
		conf,
		nil,
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
	if sb.currrentCommand == nil && !message.IsCommand() {
		interaction.Reply(sb.tgBotAPI, message, "Hello! Start by choosing a command.")
		return
	}

	if message.IsCommand() {
		sb.handleCommand(message)
		return
	}

	sb.currrentCommand.Handle(message)
}

func (sb *StickiesBot) handleCommand(message *tgbotapi.Message) {
	switch message.Command() {
	case "add":
		sb.currrentCommand = command.NewAddStickerCommand(sb.tgBotAPI, sb.conf.WalnutStickerSetName, sb.conf.WalnutStickerUserID)
		break
	default:
		interaction.Reply(sb.tgBotAPI, message, "Sorry, I don't understand that command. Please pick a command from the list.")
		sb.currrentCommand = nil
	}

	if sb.currrentCommand != nil {
		sb.currrentCommand.Start(message)
	}
}
