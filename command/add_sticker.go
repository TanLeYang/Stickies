package command

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func AddSticker(
	botapi *tgbotapi.BotAPI,
	message *tgbotapi.Message,
	stickerSetName string,
	stickerSetUserID int64) {

	sticker := message.Sticker
	if sticker == nil {
		return
	}

	_, err := addWithExistingSticker(botapi, sticker, stickerSetName, stickerSetUserID, "😡")
	if err != nil {
		log.Printf("Failed to add existing sticker to sticker set, err: %s\n", err)
	}
}

func addWithExistingSticker(
	botapi *tgbotapi.BotAPI,
	sticker *tgbotapi.Sticker,
	stickerSetName string,
	stickerSetUserID int64,
	emoji string) (*tgbotapi.APIResponse, error) {

	fileID := tgbotapi.FileID(sticker.FileID)
	config := tgbotapi.AddStickerConfig{
		UserID:       stickerSetUserID,
		Name:         stickerSetName,
		PNGSticker:   fileID,
		TGSSticker:   nil,
		Emojis:       emoji,
		MaskPosition: nil,
	}
	return botapi.Request(config)
}
