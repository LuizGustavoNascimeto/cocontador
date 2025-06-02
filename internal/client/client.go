package client

import (
	"context"
	"fmt"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store"
	waLog "go.mau.fi/whatsmeow/util/log"

	"cocontador/internal/config"
)

// MyClient é o “wrapper” que contém o whatsmeow.Client e informações de configuração
type MyClient struct {
	WAClient       *whatsmeow.Client
	eventHandlerID uint32
	cfg            *config.Config
}

// NewWhatsMeowClient cria o *whatsmeow.Client a partir do store e do logger
func NewWhatsMeowClient(deviceStore *store.Device, logger waLog.Logger) *whatsmeow.Client {
	return whatsmeow.NewClient(deviceStore, logger)
}

// NewMyClient recebe o whatsmeow.Client já criado e a configuração
func NewMyClient(wac *whatsmeow.Client, cfg *config.Config) *MyClient {
	return &MyClient{
		WAClient: wac,
		cfg:      cfg,
	}
}

func (mc *MyClient) RegisterHandler() {
	mc.eventHandlerID = mc.WAClient.AddEventHandler(mc.eventHandler)
}

func (mc *MyClient) UnregisterHandler() {
	mc.WAClient.RemoveEventHandler(mc.eventHandlerID)
}

// Connect estabelece a conexão com o WhatsApp, mostrando QR se for a primeira vez
func (mc *MyClient) Connect() error {
	// Se não houver sessão, mostra QR Code
	if mc.WAClient.Store.ID == nil {
		qrChan, _ := mc.WAClient.GetQRChannel(context.Background())
		if err := mc.WAClient.Connect(); err != nil {
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
		// Sessão salva
		if err := mc.WAClient.Connect(); err != nil {
			return err
		}
	}
	return nil
}

// Disconnect fecha a conexão com o WhatsApp
func (mc *MyClient) Disconnect() {
	mc.WAClient.Disconnect()
}
