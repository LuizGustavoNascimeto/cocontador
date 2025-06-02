package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"cocontador/internal/client"
	"cocontador/internal/config"
	store "cocontador/internal/database"

	waLog "go.mau.fi/whatsmeow/util/log"
)

func main() {
	// Carrega configurações (variáveis de ambiente, .env etc.)
	cfg, err := config.LoadConfig(".env")
	if err != nil {
		log.Fatalf("erro ao carregar configuração: %v", err)
	}

	// Inicializa o log do SQLite (nível WARN)
	dbLogger := waLog.Stdout("Database", "WARN", true)

	// Inicializa o store (SQLite) e obtém deviceStore
	deviceStore, err := store.NewDeviceStore(context.Background(), cfg.DatabaseURL, dbLogger)
	if err != nil {
		log.Fatalf("erro ao inicializar store: %v", err)
	}

	// Cria o client do WhatsMeow
	clientLogger := waLog.Stdout("Client", "WARN", true)
	waClient := client.NewWhatsMeowClient(deviceStore, clientLogger)

	// Cria instância de MyClient (wrapper)
	myClient := client.NewMyClient(waClient, cfg)

	// Registra event handler
	myClient.RegisterHandler()

	// Conecta ao WhatsApp
	if err := myClient.Connect(); err != nil {
		log.Fatalf("erro ao conectar ao WhatsApp: %v", err)
	}

	// Daqui em diante, apenas aguarda sinal de interrupção para desconectar
	fmt.Println("⏳ Aguardando conexão estabilizar...")
	time.Sleep(2 * time.Second)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	<-sigCh

	fmt.Println("🔌 Desconectando...")
	myClient.UnregisterHandler()
	myClient.Disconnect()
}
