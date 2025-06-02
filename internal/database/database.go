package database

import (
	"context"

	"github.com/mattn/go-sqlite3"
	"go.mau.fi/whatsmeow/store"
	"go.mau.fi/whatsmeow/store/sqlstore"
	waLog "go.mau.fi/whatsmeow/util/log"
)

func NewDeviceStore(ctx context.Context, dbURL string, logger waLog.Logger) (*store.Device, error) {
	sqlite3.Version()

	container, err := sqlstore.New(ctx, "sqlite3", "file:examplestore.db?_foreign_keys=on", logger)
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
