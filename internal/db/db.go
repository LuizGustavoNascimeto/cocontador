package db

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"

	"cocontador/internal/config"
	"log"

	_ "github.com/lib/pq"
)

var (
	db        *sql.DB
	once      sync.Once
	closeOnce sync.Once
	err       error
)

// GetDB retorna a instância única do pool de conexões PostgreSQL
func GetDB() (*sql.DB, error) {
	once.Do(func() {
		db, err = initDB()
	})
	return db, err
}

// initDB inicializa a conexão com o PostgreSQL e configura o pool
func initDB() (*sql.DB, error) {
	cfg := config.Load()

	// Abre a conexão com o banco
	database, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("erro ao abrir banco de dados: %w", err)
	}

	// Configura o pool de conexões
	database.SetMaxOpenConns(25)                 // Máximo de conexões abertas
	database.SetMaxIdleConns(5)                  // Máximo de conexões inativas
	database.SetConnMaxLifetime(5 * time.Minute) // Tempo máximo de vida da conexão

	// Valida a conexão com o banco
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := database.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("erro ao conectar ao banco de dados: %w", err)
	}

	log.Println("Conexão com PostgreSQL estabelecida com sucesso")
	return database, nil
}

// Close fecha a conexão com o banco de dados
func Close() error {
	var closeErr error
	closeOnce.Do(func() {
		if db != nil {
			if err := db.Close(); err != nil {
				closeErr = fmt.Errorf("erro ao fechar conexão: %w", err)
				return
			}
			log.Println("Conexão com PostgreSQL fechada")
		}
	})

	return closeErr
}

// Deprecated: Use GetDB() em vez disso
func Connect() (*sql.DB, error) {
	return GetDB()
}
