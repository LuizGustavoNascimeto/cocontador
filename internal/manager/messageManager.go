package manager

import (
	"cocontador/internal/user"
	"fmt"

	"go.mau.fi/whatsmeow/types/events"
)

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
	user_handler, err := user.NewHandler()
	if err != nil {
		return fmt.Errorf("erro ao inicializar handler de usuario: %w", err)
	}
	users_map, err := user_handler.ListAll()
	if err != nil {
		return fmt.Errorf("erro ao listar usuarios: %w", err)
	}
	u := users_map[msg.Info.Sender.String()]
	if u.ID == "" {
		fmt.Printf("Usuario %s nao encontrado, criando novo...\n", msg.Info.Sender.String())
		new_user := user.User{
			ID:   msg.Info.Sender.String(),
			Name: msg.Info.PushName,
		}
		new_user, err := user_handler.Create(new_user)
		if err != nil {
			return fmt.Errorf("erro ao criar usuario: %w", err)
		}

	}

	text := extractMessageText(msg)
	if text == "" {
		return nil
	}
	fmt.Println("Received a message!", text)
	// Aqui você pode chamar saveMessage(pgDB, ...) quando quiser
	return nil
}
