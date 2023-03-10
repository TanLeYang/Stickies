package main

import (
	"bufio"
	"context"
	"log"
	"os"

	"github.com/TanLeYang/stickies/command"
	"github.com/TanLeYang/stickies/interaction"
	stickiesset "github.com/TanLeYang/stickies/stickies_set"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type StickiesBot struct {
	tgBotAPI         *tgbotapi.BotAPI
	stickiesSetRepo  stickiesset.StickiesSetRepository
	chatToHandlerMap map[int64]*UpdateHandler
	botName          string
}

func NewStickiesBot(tgBotAPI *tgbotapi.BotAPI, stickiesSetRepo stickiesset.StickiesSetRepository, botName string) *StickiesBot {
	return &StickiesBot{
		tgBotAPI:         tgBotAPI,
		stickiesSetRepo:  stickiesSetRepo,
		chatToHandlerMap: map[int64]*UpdateHandler{},
		botName:          botName,
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
			handler := sb.getHandlerForChat(update.FromChat().ID)
			go handler.handleUpdate(update)
		}
	}
}

func (sb *StickiesBot) getHandlerForChat(chatid int64) *UpdateHandler {
	handler, ok := sb.chatToHandlerMap[chatid]
	if ok {
		return handler
	} else {
		newHandler := UpdateHandler{
			tgBotAPI:        sb.tgBotAPI,
			botName:         sb.botName,
			stickiesSetRepo: sb.stickiesSetRepo,
		}
		sb.chatToHandlerMap[chatid] = &newHandler

		return &newHandler
	}
}

type UpdateHandler struct {
	tgBotAPI        *tgbotapi.BotAPI
	botName         string
	currentCommand  command.Command
	stickiesSetRepo stickiesset.StickiesSetRepository
}

func (h *UpdateHandler) handleUpdate(update tgbotapi.Update) {
	switch {
	// Handle messages
	case update.Message != nil:
		h.handleMessage(update.Message)
		break

	// Handle button clicks
	case update.CallbackQuery != nil:
		log.Fatalf("Handling button clicks not implemented yet")
	}
}

func (h *UpdateHandler) handleMessage(message *tgbotapi.Message) {
	if h.currentCommand == nil && !message.IsCommand() {
		interaction.Reply(h.tgBotAPI, message, "Hello! Start by choosing a command.")
		return
	}

	if message.IsCommand() {
		h.handleCommand(message)
		return
	}

	if h.currentCommand.Handle(message) == command.CommandComplete {
		h.currentCommand = nil
	}
}

func (h *UpdateHandler) handleCommand(message *tgbotapi.Message) {
	switch message.Command() {
	case "addsticker":
		h.currentCommand = command.NewAddStickerCommand(h.tgBotAPI, h.stickiesSetRepo)
		break
	case "createset":
		h.currentCommand = command.NewCreateStickiesSetCommand(h.tgBotAPI, h.botName, h.stickiesSetRepo)
		break
	case "done":
		interaction.Reply(h.tgBotAPI, message, "Ok!")
		h.currentCommand = nil
		break
	case "list":
		h.currentCommand = command.NewListCommand(h.tgBotAPI, h.stickiesSetRepo)
	default:
		interaction.Reply(h.tgBotAPI, message, "Sorry, I don't understand that command. Please pick a command from the list.")
		h.currentCommand = nil
	}

	if h.currentCommand != nil {
		if h.currentCommand.Start(message) == command.CommandComplete {
			h.currentCommand = nil
		}
	}
}
