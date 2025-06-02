package client

import (
    "context"
    "fmt"
    "os"

    "go.mau.fi/whatsmeow/types/events"
    "go.mau.fi/whatsmeow/proto/waE2E"
    "google.golang.org/protobuf/proto"
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
        if messageText == targetEmoji {
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
    timestamp := msg.Info.Timestamp.Format("15:04:05")
    remetente := msg.Info.Sender.User
    pushName := msg.Info.PushName
    chat := msg.Info.Chat

    fmt.Printf("✅ Emoji válido recebido de %s às %s\n", remetente, timestamp)
    resposta := fmt.Sprintf(
        "Mensagem recebida de %s às %s. Emoji processado com sucesso!\nAss: Cocontador",
        pushName, timestamp,
    )

    // Envia mensagem de volta para o grupo ou contato
    if _, err := mc.WAClient.SendMessage(context.Background(), chat, &waE2E.Message{
        Conversation: proto.String(resposta),
    }); err != nil {
        fmt.Fprintf(os.Stderr, "erro ao enviar resposta: %v\n", err)
    } else {
        fmt.Println("✅ Resposta enviada!")
    }
}
