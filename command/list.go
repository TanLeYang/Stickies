package command

import (
	"fmt"
	"log"

	"github.com/TanLeYang/stickies/interaction"
	stickiesset "github.com/TanLeYang/stickies/stickies_set"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type List struct {
	botapi          *tgbotapi.BotAPI
	stickiesSetRepo stickiesset.StickiesSetRepository
}

func NewListCommand(botapi *tgbotapi.BotAPI, stickiesSetRepo stickiesset.StickiesSetRepository) *List {
	return &List{
		botapi,
		stickiesSetRepo,
	}
}

func (c *List) Start(message *tgbotapi.Message) bool {
	userID := message.From.ID
	setsOwnedByUser, err := c.stickiesSetRepo.GetByOwner(userID)
	if err != nil {
		interaction.GenericErrorReply(c.botapi, message)
		return false
	}

	for _, set := range setsOwnedByUser {
		setName := set.TgStickerSetName
		uniqueCode := set.UniqueCode
		interaction.Reply(c.botapi, message, fmt.Sprintf("Set: %s", setName))
		interaction.Reply(c.botapi, message, uniqueCode)
	}

	return false
}

func (c *List) Handle(message *tgbotapi.Message) bool {
	log.Printf("List command does not expect further user interaction, message received: %v", message)
	return false
}
