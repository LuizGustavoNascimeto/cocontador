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
	_ "github.com/mattn/go-sqlite3"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"
	"google.golang.org/protobuf/proto"
)

type MyClient struct {
	WAClient       *whatsmeow.Client
	eventHandlerID uint32
}

// NewMyClient cria uma nova instância do MyClient
func NewMyClient(client *whatsmeow.Client) *MyClient {
	return &MyClient{
		WAClient: client,
	}
}

// register registra o event handler
func (mycli *MyClient) register() {
	mycli.eventHandlerID = mycli.WAClient.AddEventHandler(mycli.eventHandler)
}

// unregister remove o event handler
func (mycli *MyClient) unregister() {
	mycli.WAClient.RemoveEventHandler(mycli.eventHandlerID)
}

// connect conecta ao WhatsApp
func (mycli *MyClient) connect() error {
	if mycli.WAClient.Store.ID == nil {
		// Primeiro uso - precisa do QR code
		qrChan, _ := mycli.WAClient.GetQRChannel(context.Background())
		err := mycli.WAClient.Connect()
		if err != nil {
			return err
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
		err := mycli.WAClient.Connect()
		if err != nil {
			return err
		}
	}
	return nil
}

// disconnect desconecta do WhatsApp
func (mycli *MyClient) disconnect() {
	mycli.WAClient.Disconnect()
}

func (mycli *MyClient) eventHandler(evt interface{}) {
	var GROUP_JID string = os.Getenv("GROUP_JID")
	var TARGET_EMOJI string = os.Getenv("EMOJI")

	switch v := evt.(type) {
	case *events.Message:
		// Filtra apenas mensagens do grupo específico
		if v.Info.Chat.String() != GROUP_JID {
			return
		}

		// Log para debug - ver todas as mensagens recebidas
🧿
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
			mycli.processValidEmoji(v)
		}

	case *events.Connected:
		fmt.Println("🔗 Conectado ao WhatsApp")
	case *events.Disconnected:
		fmt.Println("❌ Desconectado do WhatsApp")
	}
}

func (mycli *MyClient) processValidEmoji(msg *events.Message) {
	fmt.Printf("✅ Processando emoji válido de: %s às %s\n",
		msg.Info.Sender.User,
		msg.Info.Timestamp.Format("15:04:05"))
	// Sua lógica de negócio aqui
	conversationMessage := "Mensagem recebida de " + msg.Info.PushName + " às " + msg.Info.Timestamp.Format("15:04:05") + ". Emoji processado com sucesso!\n Ass: Cocontador"
	mycli.WAClient.SendMessage(context.Background(), msg.Info.Chat, &waE2E.Message{Conversation: proto.String(conversationMessage)})
	fmt.Println("✅ Resposta enviada!")

}

func main() {
	// Carregar as variáveis de ambiente
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

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

	// Cria instância do MyClient
	myClient := NewMyClient(client)

	// Registra o event handler
	myClient.register()

	// Conecta ao WhatsApp
	err = myClient.connect()
	if err != nil {
		panic(err)
	}
	//exibir todos os grupos

	// Aguarda conexão estabilizar
	fmt.Println("⏳ Aguardando conexão...")
	time.Sleep(2 * time.Second)

	// Mantém o programa rodando
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	fmt.Println("Desconectando...")
	myClient.unregister()
	myClient.disconnect()
}
