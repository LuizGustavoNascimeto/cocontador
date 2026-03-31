package main

import (
	"cocontador/internal/db"
	whatsmeow "cocontador/internal/middleware"
	"log"
)

func main() {
	// Inicializa o pool de conexões PostgreSQL (singleton)
	conn, err := db.GetDB()
	if err != nil {
		log.Fatalf("Erro ao conectar ao banco de dados: %v", err)
	}
	defer db.Close()

	// Testa a conexão
	if err := conn.Ping(); err != nil {
		log.Fatalf("Erro ao fazer ping no banco: %v", err)
	}
	log.Println("✓ Pool de conexões PostgreSQL inicializado com sucesso")

	whatsmeow.StartWhatsmeow()
	log.Println("✓ Desligando...")
}
