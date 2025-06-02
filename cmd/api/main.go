package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/mattn/go-sqlite3"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"
)

func eventHandler(evt interface{}) {
	var GROUP_JID string = os.Getenv("GROUP_JID")
	var TARGET_EMOJI string = os.Getenv("EMOJI")

	switch v := evt.(type) {
	case *events.Message:
		// Filtra apenas mensagens do grupo específico
		if v.Info.Chat.String() != GROUP_JID {
			return
		}

		// Log para debug - ver todas as mensagens recebidas
		fmt.Printf("📱 Mensagem recebida de: %s\n", v.Info.Sender.User)

		// Obtém o texto da mensagem de forma mais abrangente
		var messageText string

		// Verifica diferentes tipos de mensagem
		if v.Message.GetConversation() != "" {
			messageText = v.Message.GetConversation()
			fmt.Printf("Texto (Conversation): '%s'\n", messageText)
		} else if v.Message.GetExtendedTextMessage() != nil {
			messageText = v.Message.GetExtendedTextMessage().GetText()
			fmt.Printf("Texto (Extended): '%s'\n", messageText)
		} else {
			// Mensagem sem texto identificável
			fmt.Println("Mensagem sem texto ou tipo não suportado")
			return
		}

		// Verifica se é igual ao emoji
		if messageText == TARGET_EMOJI {
			fmt.Printf("Emoji encontrado! Remetente: %s\n", v.Info.Sender.User)
			processValidEmoji(v)
		}

	case *events.Connected:
		fmt.Println("🔗 Conectado ao WhatsApp")
	case *events.Disconnected:
		fmt.Println("❌ Desconectado do WhatsApp")
	}
}

func processValidEmoji(msg *events.Message) {
	fmt.Printf(" Processando emoji válido de: %s às %s\n",
		msg.Info.Sender.User,
		msg.Info.Timestamp.Format("15:04:05"))

	// Sua lógica de negócio aqui
}

func main() {
	// Carregar as variáveis de ambiente
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

	sqlite3.Version() // Verifica se o driver está funcionando
	dbLog := waLog.Stdout("Database", "WARN", true)
	ctx := context.Background()

	container, err := sqlstore.New(ctx, "sqlite3", "file:examplestore.db?_foreign_keys=on", dbLog)
	if err != nil {
		panic(err)
	}

	deviceStore, err := container.GetFirstDevice(ctx)
	if err != nil {
		panic(err)
	}

	clientLog := waLog.Stdout("Client", "WARN", true)
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
	fmt.Println("⏳ Aguardando conexão...")
	time.Sleep(2 * time.Second)

	// Mantém o programa rodando
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	fmt.Println("Desconectando...")
	client.Disconnect()
}
