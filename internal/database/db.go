package database

import (
	"fmt"
	"money_map/internal/modules/account"
	"money_map/internal/modules/bank"
	"money_map/internal/modules/card"
	"money_map/internal/modules/category"
	"money_map/internal/modules/tags"
	"money_map/internal/modules/transaction"
	"money_map/internal/modules/user"
	"sync"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	dbConn *gorm.DB
	once   sync.Once
)

func Connect(dns string) (*gorm.DB, error) {
	var initError error
	once.Do(func() {
		var err error
		dbConn, err = gorm.Open(postgres.Open(dns), &gorm.Config{})
		if err != nil {
			initError = fmt.Errorf("failed to open database: %v", err)
			return
		}
	})
	if initError != nil {
		return nil, initError
	}
	return dbConn, nil
}

func Close() error {
	sqlDB, err := dbConn.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB from gorm DB: %v", err)
	}
	if err := sqlDB.Close(); err != nil {
		return fmt.Errorf("failed to close database connection: %v", err)
	}
	return nil
}

func Migrate() error {
	if dbConn == nil {
		return fmt.Errorf("database connection is not initialized")
	}
	err := dbConn.AutoMigrate(
		&bank.Bank{},
		&user.User{},
		&category.Category{},
		&transaction.Transaction{},
		&card.Card{},
		&tags.Tags{},
		&account.Account{},
	)
	if err != nil {
		return fmt.Errorf("failed to migrate database: %v", err)
	}
	return nil
}
