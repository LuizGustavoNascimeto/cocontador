package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL string
	GroupJID    string
	TargetEmoji string
}

func LoadConfig(envFile string) (*Config, error) {
	// Tenta carregar .env (se existir)
	if err := godotenv.Load(envFile); err != nil {
		// Não fatal se o arquivo não existir; apenas loga
		fmt.Println("⚠️ Não foi possível carregar .env ou arquivo não existe. Usando variáveis de ambiente do sistema.")
	}

	group := os.Getenv("GROUP_JID")
	if group == "" {
		return nil, fmt.Errorf("GROUP_JID não definido no ambiente")
	}

	emoji := os.Getenv("EMOJI")
	if emoji == "" {
		return nil, fmt.Errorf("EMOJI não definido no ambiente")
	}
	postgresURL := os.Getenv("DB_URL")
	if postgresURL == "" {
		postgresURL = "postgres://postgres:@localhost:5432/whatsapp_db?sslmode=disable"
	}

	return &Config{
		GroupJID:    group,
		TargetEmoji: emoji,
		DatabaseURL: postgresURL,
	}, nil
}
