package manager

import (
	"cocontador/internal/barrigade"
	"cocontador/internal/user"
	"cocontador/internal/util"
	"fmt"

	"go.mau.fi/whatsmeow/types/events"
)

var (
	user_handler      *user.UserHandler
	barrigade_handler *barrigade.BarrigadeHandler
)

func init() {
	var err error
	user_handler, err = user.NewHandler()
	if err != nil {
		panic(fmt.Sprintf("erro ao criar user handler: %v", err))
	}
	barrigade_handler, err = barrigade.NewHandler()
	if err != nil {
		panic(fmt.Sprintf("erro ao criar barrigade handler: %v", err))
	}
}

func extractMessageText(msg *events.Message) string {
	text := msg.Message.GetConversation()
	if text != "" {
		return text
	}

	extended := msg.Message.GetExtendedTextMessage()
	if extended == nil {
		return ""
	}

	return extended.GetText()
}

func MesssageManager(msg *events.Message) error {
	if !MsgIsValid(msg) {
		return nil
	}
	text := extractMessageText(msg)
	if text == "" {
		return nil
	}
	fmt.Println("Received a message!", text)

	user, err := user_handler.CreateOrGet(user.User{
		ID: msg.Info.Sender.User,
	})
	if err != nil {
		return fmt.Errorf("erro ao criar ou obter usuario: %w", err)
	}

	_, err = barrigade_handler.Create(barrigade.Barrigade{
		User_id:    user.ID,
		Created_at: msg.Info.Timestamp,
	})
	if err != nil {
		return fmt.Errorf("erro ao criar barrigade: %w", err)
	}

	return nil
}
func MsgIsValid(msg *events.Message) bool {
	text := extractMessageText(msg)
	hasEmoji, _ := util.Contains(text, "👍")
	if hasEmoji {
		return true
	}
	return false
}
