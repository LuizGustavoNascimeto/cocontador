package client

import (
	utils "cocontador/pkg"
	"context"
	"fmt"
	"os"

	"go.mau.fi/whatsmeow/types/events"
)

// eventHandler é o callback que o whatsmeow chama para cada evento
func (mc *MyClient) eventHandler(evt interface{}) {
	groupJID := mc.cfg.GroupJID
	targetEmoji := mc.cfg.TargetEmoji

	switch e := evt.(type) {
	case *events.Message:
		// Só processa mensagens do grupo configurado
		if e.Info.Chat.String() != groupJID {
			return
		}

		// Extrai o texto da mensagem (conv + extended)
		var messageText string
		if conv := e.Message.GetConversation(); conv != "" {
			messageText = conv
		} else if ext := e.Message.GetExtendedTextMessage(); ext != nil {
			messageText = ext.GetText()
		} else {
			// Sem texto ou tipo não suportado
			return
		}

		// Se for o emoji certo, processa
		res, _ := utils.Contains(messageText, targetEmoji)
		if res {
			mc.processValidEmoji(e)
		}

	case *events.Connected:
		fmt.Println("🔗 Conectado ao WhatsApp")
	case *events.Disconnected:
		fmt.Println("❌ Desconectado do WhatsApp")
	}
}

// processValidEmoji contém a lógica de negócio (respondendo, logando, etc.)
func (mc *MyClient) processValidEmoji(msg *events.Message) {
	timestamp := getCorrectTime(msg)
	remetente := msg.Info.Sender.User
	// pushName := msg.Info.PushName
	chat := msg.Info.Chat
	senderJID := msg.Info.Sender
	targetMessageID := msg.Info.ID

	//prepara a resposta
	fmt.Printf("✅ Emoji válido recebido de %s às %s\n", remetente, timestamp)

	// reage a mensagem
	_, err := mc.WAClient.SendMessage(context.Background(), chat, mc.WAClient.BuildReaction(chat, senderJID, targetMessageID, "👍"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "erro ao enviar resposta: %v\n", err)
	} else {
		fmt.Println("✅ Resposta enviada!")
	}
}

func getCorrectTime(msg *events.Message) string {
	msgConversation := msg.Message.GetConversation()

	res, pos := utils.Contains(msgConversation, "⏰")
	if res {
		return msgConversation[pos+3:]
	}
	return msg.Info.Timestamp.Format("15:04:05") // Formata o timestamp padrão se não encontrar o emoji

}
