package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"
)

const (
	TARGET_GROUP_NAME = "Católicos suaves  ✨✨🙏" // Nome do grupo
	TARGET_EMOJI      = "💩"                     // Emoji que você quer validar
)

var TARGET_GROUP_ID string // Será preenchido dinamicamente

func eventHandler(evt interface{}) {
	switch v := evt.(type) {
	case *events.Message:
		// Filtra apenas mensagens do grupo específico
		if v.Info.Chat.String() != TARGET_GROUP_ID {
			return
		}

		// Verifica se a mensagem contém texto
		if v.Message.GetConversation() == "" &&
			v.Message.GetExtendedTextMessage() == nil {
			return
		}

		// Obtém o texto da mensagem
		var messageText string
		if v.Message.GetConversation() != "" {
			messageText = v.Message.GetConversation()
		} else if v.Message.GetExtendedTextMessage() != nil {
			messageText = v.Message.GetExtendedTextMessage().GetText()
		}
		// verifica se é igual ao emoji
		if strings.TrimSpace(messageText) == TARGET_EMOJI {
			fmt.Printf("✅ Emoji encontrado! Remetente: %s\n", v.Info.Sender.User)

			// Aqui você pode aplicar sua lógica de negócio
			processValidEmoji(v)
		}

	case *events.Connected:
		fmt.Println("🔗 Conectado ao WhatsApp")

	case *events.Disconnected:
		fmt.Println("❌ Desconectado do WhatsApp")
	}
}

func processValidEmoji(msg *events.Message) {
	// Sua lógica de negócio aqui
	fmt.Printf("Processando emoji válido de: %s às %s\n",
		msg.Info.Sender.User,
		msg.Info.Timestamp.Format("15:04:05"))

	// Exemplo: contar, salvar em BD, enviar notificação, etc.
}

func findGroupByName(client *whatsmeow.Client, groupName string) (string, error) {
	groups, err := client.GetJoinedGroups()
	if err != nil {
		return "", err
	}

	for _, group := range groups {
		info, err := client.GetGroupInfo(group)
		if err != nil {
			continue 5
		}
		if info.Name == groupName {
			return group.String(), nil
		}
	}
	return "", fmt.Errorf("grupo '%s' não encontrado", groupName)
}
func main() {
	dbLog := waLog.Stdout("Database", "DEBUG", true)
	ctx := context.Background()

	container, err := sqlstore.New(ctx, "sqlite3", "file:examplestore.db?_foreign_keys=on", dbLog)
	if err != nil {
		panic(err)
	}

	deviceStore, err := container.GetFirstDevice(ctx)
	if err != nil {
		panic(err)
	}

	clientLog := waLog.Stdout("Client", "DEBUG", true)
	client := whatsmeow.NewClient(deviceStore, clientLog)
	client.AddEventHandler(eventHandler)

	// Conecta ao WhatsApp
	if client.Store.ID == nil {
		// Primeiro uso - precisa do QR code
		qrChan, _ := client.GetQRChannel(context.Background())
		err = client.Connect()
		if err != nil {
			panic(err)
		}
		for evt := range qrChan {
			if evt.Event == "code" {
				fmt.Println("QR code:", evt.Code)
			} else {
				fmt.Println("Login event:", evt.Event)
			}
		}
	} else {
		// Já tem sessão salva
		err = client.Connect()
		if err != nil {
			panic(err)
		}
	}

	// Aguarda conexão estabilizar
	fmt.Println("Aguardando conexão...")
	time.Sleep(3 * time.Second)

	// Busca o ID do grupo pelo nome
	groupID, err := findGroupByName(client, TARGET_GROUP_NAME)
	if err != nil {
		panic(fmt.Sprintf("Erro ao buscar grupo: %v", err))
	}

	TARGET_GROUP_ID = groupID
	fmt.Printf("✅ Grupo '%s' encontrado: %s\n", TARGET_GROUP_NAME, TARGET_GROUP_ID)

	// Mantém o programa rodando
	select {}
}
