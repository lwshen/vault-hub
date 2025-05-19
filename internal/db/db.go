package db

import (
	"fmt"
	"log/slog"

	"github.com/lwshen/vault-hub/internal/config"
	slogGorm "github.com/orandin/slog-gorm"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Open(logger *slog.Logger) error {
	gormConfig := &gorm.Config{
		Logger: slogGorm.New(slogGorm.WithHandler(logger.Handler())),
	}

	switch config.DatabaseType {
	case config.DatabaseTypeSQLite:
		db, err := OpenSQLite(gormConfig)
		if err != nil {
			return err
		}
		DB = db
	}

	if DB == nil {
		return fmt.Errorf("database not initialized")
	}

	if err := checkConnection(); err != nil {
		return err
	}

	return nil
}

func OpenSQLite(gormConfig *gorm.Config) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(config.DatabaseUrl), gormConfig)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func checkConnection() error {
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}

	return sqlDB.Ping()
}
