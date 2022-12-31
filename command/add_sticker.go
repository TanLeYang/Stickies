package command

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	Upload      CommandStage = "UPLOAD"
	ChooseEmoji CommandStage = "CHOOSE_EMOJI"
)

var addStickerStages = []CommandStage{
	Upload,
	ChooseEmoji,
}

type AddSticker struct {
	botapi           *tgbotapi.BotAPI
	stickerSetName   string
	stickerSetUserID int64
	currentStage     CommandStage
	stickerToAdd     tgbotapi.RequestFileData
}

func NewAddStickerCommand(botapi *tgbotapi.BotAPI, stickerSetName string, stickerSetUserID int64) *AddSticker {
	return &AddSticker{
		botapi:           botapi,
		stickerSetName:   stickerSetName,
		stickerSetUserID: stickerSetUserID,
		currentStage:     addStickerStages[0],
		stickerToAdd:     nil,
	}
}

func (c *AddSticker) Start(message *tgbotapi.Message) {
	reply(c.botapi, message, "Please send me a sticker.")
}

func (c *AddSticker) Handle(message *tgbotapi.Message) {
	switch c.currentStage {
	case Upload:
		c.uploadStage(message)
		break
	case ChooseEmoji:
		c.chooseEmojiStage(message)
		break
	default:
		log.Panicf("Undefined stage for AddSticker command: %s", c.currentStage)
	}
}

func (c *AddSticker) uploadStage(message *tgbotapi.Message) {
	sticker := message.Sticker
	if sticker == nil {
		reply(c.botapi, message, "Please include a sticker or PNG file.")
		return
	}

	c.stickerToAdd = tgbotapi.FileID(sticker.FileID)
	reply(c.botapi, message, "Thanks! Now send me an emoji that corresponds to the sticker.")
	c.currentStage = nextStage(addStickerStages, c.currentStage, true)
}

func (c *AddSticker) chooseEmojiStage(message *tgbotapi.Message) {
	emoji := message.Text

	if c.stickerToAdd == nil {
		log.Panicf("Sticker to add is nil at %s stage", c.currentStage)
	}
	_, err := c.addStickerRequest(c.stickerToAdd, emoji)
	if err != nil {
		reply(c.botapi, message, "Sorry, something went wrong, please restart and try again.")
	} else {
		reply(c.botapi, message, "Nice! Sticker has been added to the set. To add another one, send me the next sticker.")
	}

	c.currentStage = nextStage(addStickerStages, c.currentStage, true)
	c.stickerToAdd = nil
}

func (c *AddSticker) addStickerRequest(stickerToAdd tgbotapi.RequestFileData, emoji string) (*tgbotapi.APIResponse, error) {
	config := tgbotapi.AddStickerConfig{
		UserID:       c.stickerSetUserID,
		Name:         c.stickerSetName,
		PNGSticker:   stickerToAdd,
		TGSSticker:   nil,
		Emojis:       emoji,
		MaskPosition: nil,
	}

	return c.botapi.Request(config)
}
