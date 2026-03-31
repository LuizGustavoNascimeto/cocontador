package whatsmeow

import (
	"cocontador/internal/manager"
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/lib/pq"
	"github.com/mattn/go-sqlite3"
	_ "github.com/mattn/go-sqlite3"
	"github.com/mdp/qrterminal/v3"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"
)

func eventHandler(evt interface{}) {
	switch v := evt.(type) {
	case *events.Message:
		err := manager.MesssageManager(v)
		if err != nil {
			fmt.Println("Error handling message:", err)
		}
	}
}

func StartWhatsmeow() {
	ctx := context.Background()

	sqlite3.Version() // Força a inicialização do driver SQLite3

	client, err := newClient(ctx)
	if err != nil {
		panic(err)
	}

	if err := connectClient(ctx, client); err != nil {
		panic(err)
	}

	waitForShutdownSignal()
	client.Disconnect()
}

func newClient(ctx context.Context) (*whatsmeow.Client, error) {
	dbLog := waLog.Stdout("Database", "ERROR", true)
	container, err := sqlstore.New(ctx, "sqlite3", "file:examplestore.db?_foreign_keys=on", dbLog)
	if err != nil {
		return nil, err
	}

	deviceStore, err := container.GetFirstDevice(ctx)
	if err != nil {
		return nil, err
	}

	clientLog := waLog.Stdout("Client", "ERROR", true)
	client := whatsmeow.NewClient(deviceStore, clientLog)
	client.AddEventHandler(eventHandler)

	return client, nil
}

func connectClient(ctx context.Context, client *whatsmeow.Client) error {
	if client.Store.ID == nil {
		return connectWithQRCode(ctx, client)
	}

	return client.Connect()
}

func connectWithQRCode(ctx context.Context, client *whatsmeow.Client) error {
	qrChan, err := client.GetQRChannel(ctx)
	if err != nil {
		return err
	}

	if err := client.Connect(); err != nil {
		return err
	}

	for evt := range qrChan {
		if evt.Event == "code" {
			qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, os.Stdout)
			fmt.Println("QR code:", evt.Code)
			continue
		}

		fmt.Println("Login event:", evt.Event)
	}

	return nil
}

func waitForShutdownSignal() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
}
