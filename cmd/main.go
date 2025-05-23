package main

import (
	"cocontador/internal/routers"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	// Carregar variáveis de ambiente do arquivo .env
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Erro ao carregar o arquivo .env: %v", err)
	}

	// Configurar rotas
	router := routers.SetupRouter()

	// Obter a porta da API a partir das variáveis de ambiente
	port := os.Getenv("API_PORT")
	if port == "" {
		log.Fatal("A variável de ambiente API_PORT não está definida")
	}

	// Iniciar o servidor
	router.Run(":" + port)
}

// Welcome Message
