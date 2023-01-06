package command

import (
	"fmt"
	"log"

	"github.com/TanLeYang/stickies/interaction"
	stickiesset "github.com/TanLeYang/stickies/stickies_set"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	EnterUniqueCode CommandStage = "CHOOSE_UNIQUE_CODE"
	Upload          CommandStage = "UPLOAD"
	ChooseEmoji     CommandStage = "CHOOSE_EMOJI"
)

type AddSticker struct {
	botapi             *tgbotapi.BotAPI
	stickiesSetRepo    stickiesset.StickiesSetRepository
	stickiesSetToAddTo *stickiesset.StickiesSet
	currentStage       CommandStage
	stickerToAdd       tgbotapi.RequestFileData
	emoji              string
}

func NewAddStickerCommand(botapi *tgbotapi.BotAPI, stickiesSetRepo stickiesset.StickiesSetRepository) *AddSticker {
	return &AddSticker{
		botapi:          botapi,
		stickiesSetRepo: stickiesSetRepo,
		currentStage:    EnterUniqueCode,
	}
}

func (c *AddSticker) Start(message *tgbotapi.Message) {
	interaction.Reply(c.botapi, message, "Please send me the unique sharing code of the sticker set you wish to add stickers to.")
}

func (c *AddSticker) Handle(message *tgbotapi.Message) {
	switch c.currentStage {
	case EnterUniqueCode:
		c.enterUniqueCodeStage(message)
		break
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

func (c *AddSticker) enterUniqueCodeStage(message *tgbotapi.Message) {
	uniqueCode := message.Text
	stickiesSet, err := c.stickiesSetRepo.GetByUniqueCode(uniqueCode)
	if err != nil {
		interaction.Reply(c.botapi, message, "Please send me a valid unique sharing code.")
		return
	}

	interaction.Reply(c.botapi, message, fmt.Sprintf("Adding to the %s sticker set!", stickiesSet.TgStickerSetName))
	interaction.Reply(c.botapi, message, fmt.Sprintf("Please send me the sticker to add! It can be a png file or an existing sticker."))

	c.stickiesSetToAddTo = stickiesSet
	c.currentStage = Upload
}

func (c *AddSticker) uploadStage(message *tgbotapi.Message) {
	file := uploadStickerInteraction(c.botapi, message)
	if file == nil {
		return
	}

	interaction.Reply(c.botapi, message, "Thanks! Now send me an emoji that corresponds to the sticker.")

	c.stickerToAdd = file
	c.currentStage = ChooseEmoji
}

func uploadStickerInteraction(botapi *tgbotapi.BotAPI, message *tgbotapi.Message) tgbotapi.RequestFileData {
	var file tgbotapi.RequestFileData
	if message.Sticker != nil {
		file = tgbotapi.FileID(message.Sticker.FileID)
	} else if message.Document != nil && message.Document.MimeType == "image/png" {
		file = tgbotapi.FileID(message.Document.FileID)
	} else {
		interaction.Reply(botapi, message, `Please send either a sticker or PNG file
			not exceeding 512kb in size, either width or height exactly 512px.`)
		return nil
	}

	return file
}

func (c *AddSticker) chooseEmojiStage(message *tgbotapi.Message) {
	if c.stickerToAdd == nil {
		log.Panicf("Sticker to add is nil at %s stage", c.currentStage)
	}

	emoji := chooseEmojiInteraction(c.botapi, message)
	c.emoji = emoji

	_, err := c.addStickerRequest()
	if err != nil {
		interaction.Reply(c.botapi, message, "Sorry, something went wrong, please restart and try again.")
	} else {
		interaction.Reply(c.botapi, message,
			"Nice! Sticker has been added to the set. To add another one, send me the next sticker. Use the /done command once you're done.")
	}

	c.currentStage = Upload
	c.stickiesSetToAddTo = nil
	c.stickerToAdd = nil
	c.emoji = ""
}

func chooseEmojiInteraction(botapi *tgbotapi.BotAPI, message *tgbotapi.Message) string {
	return message.Text
}

func (c *AddSticker) addStickerRequest() (*tgbotapi.APIResponse, error) {
	config := tgbotapi.AddStickerConfig{
		UserID:       c.stickiesSetToAddTo.Owner,
		Name:         c.stickiesSetToAddTo.TgStickerSetName,
		PNGSticker:   c.stickerToAdd,
		TGSSticker:   nil,
		Emojis:       c.emoji,
		MaskPosition: nil,
	}

	return c.botapi.Request(config)
}
