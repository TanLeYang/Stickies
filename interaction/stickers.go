package interaction

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

func UploadStickerInteraction(botapi *tgbotapi.BotAPI, message *tgbotapi.Message) tgbotapi.RequestFileData {
	var file tgbotapi.RequestFileData
	if message.Sticker != nil {
		file = tgbotapi.FileID(message.Sticker.FileID)
	} else if message.Document != nil && message.Document.MimeType == "image/png" {
		file = tgbotapi.FileID(message.Document.FileID)
	} else {
		Reply(botapi, message,
			"Please send either a sticker or PNG file not exceeding 512kb in size, either width or height exactly 512px.")
		return nil
	}

	return file
}

func ChooseEmojiInteraction(botapi *tgbotapi.BotAPI, message *tgbotapi.Message) string {
	return message.Text
}
