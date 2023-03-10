package command

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/TanLeYang/stickies/interaction"
	stickiesset "github.com/TanLeYang/stickies/stickies_set"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	ChooseSetName        CommandStage = "CHOOSE_SET_NAME"
	UploadInitialSticker CommandStage = "UPLOAD_INITIAL_STICKER"
	ChooseInitialEmoji   CommandStage = "CHOOSE_INITIAL_EMOJI"
)

type CreateStickiesSet struct {
	botapi          *tgbotapi.BotAPI
	botName         string
	stickiesSetRepo stickiesset.StickiesSetRepository
	currentStage    CommandStage
	stickerSetName  string
	stickerToAdd    *tgbotapi.RequestFileData
	emoji           string
}

func NewCreateStickiesSetCommand(botapi *tgbotapi.BotAPI, botName string, stickiesSetRepo stickiesset.StickiesSetRepository) *CreateStickiesSet {
	return &CreateStickiesSet{
		botapi:          botapi,
		botName:         botName,
		stickiesSetRepo: stickiesSetRepo,
		currentStage:    ChooseSetName,
	}
}

func (c *CreateStickiesSet) Start(message *tgbotapi.Message) CommandOngoingStatus {
	interaction.Reply(c.botapi, message, "Please choose a sticker set name.")
	return CommandOngoing
}

func (c *CreateStickiesSet) Handle(message *tgbotapi.Message) CommandOngoingStatus {
	switch c.currentStage {
	case ChooseSetName:
		return c.chooseSetNameStage(message)
	case UploadInitialSticker:
		return c.uploadInitialStickerStage(message)
	case ChooseInitialEmoji:
		return c.chooseInitialEmojiStage(message)
	default:
		log.Panicf("Undefined stage for CreateStickiesSet command: %s", c.currentStage)
		return CommandComplete
	}
}

func (c *CreateStickiesSet) chooseSetNameStage(message *tgbotapi.Message) CommandOngoingStatus {
	setName := message.Text
	if len(setName) == 0 {
		interaction.Reply(c.botapi, message, "Please send me plaintext!")
		return CommandOngoing
	}

	c.stickerSetName = setName
	c.currentStage = UploadInitialSticker

	interaction.Reply(c.botapi, message, "Got it! Now please send me a png file or sticker to add into the set.")

	return CommandOngoing
}

func (c *CreateStickiesSet) uploadInitialStickerStage(message *tgbotapi.Message) CommandOngoingStatus {
	file := interaction.UploadStickerInteraction(c.botapi, message)
	if file == nil {
		return CommandOngoing
	}

	c.stickerToAdd = &file
	c.currentStage = ChooseInitialEmoji

	interaction.Reply(c.botapi, message, "Ok! Now send me an emoji that corresponds to the sticker.")

	return CommandOngoing
}

func (c *CreateStickiesSet) chooseInitialEmojiStage(message *tgbotapi.Message) CommandOngoingStatus {
	emoji := interaction.ChooseEmojiInteraction(c.botapi, message)
	c.emoji = emoji

	replyGenericError := func() {
		interaction.Reply(c.botapi, message, "Sorry, something went wrong. Please try again with another name!")
	}

	tgStickerSetName := c.formatSetName(c.stickerSetName)
	_, tgErr := c.createTgStickerPack(message.From.ID, tgStickerSetName, c.stickerSetName)
	if tgErr != nil {
		replyGenericError()
		return CommandOngoing
	}

	randomIdentifer := generateRandomIdentifer(c.stickerSetName)
	persistErr := c.stickiesSetRepo.Create(stickiesset.StickiesSet{
		Owner:            message.From.ID,
		TgStickerSetName: tgStickerSetName,
		UniqueCode:       randomIdentifer,
	})
	if persistErr != nil {
		replyGenericError()
		return CommandOngoing
	}

	interaction.Reply(c.botapi, message, fmt.Sprintf(
		"Your sticker set has been created! You can find it at telegram.me/addstickers/%s. \n"+
			"Use the code that will be sent in the next message to add stickers to the set using the /addsticker command. \n"+
			"Anyone with the code can contribute!",
		tgStickerSetName,
	))
	interaction.Reply(c.botapi, message, randomIdentifer)

	return CommandComplete
}

func (c *CreateStickiesSet) createTgStickerPack(userID int64, tgStickerSetName string, stickerSetTitle string) (*tgbotapi.APIResponse, error) {
	conf := tgbotapi.NewStickerSetConfig{
		UserID:        userID,
		Name:          tgStickerSetName,
		Title:         stickerSetTitle,
		PNGSticker:    *c.stickerToAdd,
		TGSSticker:    nil,
		Emojis:        c.emoji,
		ContainsMasks: false,
		MaskPosition:  nil,
	}

	return c.botapi.Request(conf)
}

func (c *CreateStickiesSet) formatSetName(name string) string {
	return fmt.Sprintf("%s_by_%s", name, c.botName)
}

func generateRandomIdentifer(setName string) string {
	length := 10
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, length)
	rand.Read(b)
	randomPart := fmt.Sprintf("%x", b)[:length]

	return fmt.Sprintf("%s-%s", setName, randomPart)
}
