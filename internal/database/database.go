package database

import (
	"cocontador/internal/config"
	"context"
	"database/sql"

	_ "github.com/lib/pq"
	"github.com/mattn/go-sqlite3"
	"go.mau.fi/whatsmeow/store"
	"go.mau.fi/whatsmeow/store/sqlstore"
	waLog "go.mau.fi/whatsmeow/util/log"
)

func NewDeviceStore(ctx context.Context, dbURL string, logger waLog.Logger) (*store.Device, error) {
	sqlite3.Version()

	container, err := sqlstore.New(ctx, "postgres", dbURL, nil)
	// If you want multiple sessions, remember their JIDs and use .GetDevice(jid) or .GetAllDevices() instead.

	if err != nil {
		return nil, err
	}

	deviceStore, err := container.GetFirstDevice(ctx)
	if err != nil {
		return nil, err
	}
	return deviceStore, nil

}

func OpenConn() (*sql.DB, error) {
	cfg, err := config.LoadConfig(".env")
	if err != nil {
		panic(err)
	}
	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		panic(err)
	}

	err = db.Ping()

	return db, err
}
