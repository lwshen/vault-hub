package model

import (
	"fmt"
	"log/slog"

	"github.com/lwshen/vault-hub/internal/config"
	slogGorm "github.com/orandin/slog-gorm"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Open(logger *slog.Logger) error {
	gormConfig := &gorm.Config{
		Logger: slogGorm.New(slogGorm.WithHandler(logger.Handler())),
	}

	var d *gorm.DB
	var err error

	switch config.DatabaseType {
	case config.DatabaseTypeSQLite:
		d, err = openSQLite(gormConfig)
	case config.DatabaseTypeMySQL:
		d, err = openMySQL(gormConfig)
	case config.DatabaseTypePostgres:
		d, err = openPostgres(gormConfig)
	default:
		err = fmt.Errorf("unsupported database type: %s", config.DatabaseType)
	}

	if err != nil {
		return err
	}
	DB = d

	if DB == nil {
		// This case should ideally not be reached if the Open<Type> functions
		// and the default case correctly handle errors and return values.
		// However, it's a safeguard.
		return fmt.Errorf("database not initialized despite no explicit error")
	}

	if err := checkConnection(); err != nil {
		return err
	}

	if err := migrate(); err != nil {
		return err
	}

	return nil
}

func openSQLite(gormConfig *gorm.Config) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(config.DatabaseUrl), gormConfig)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func openMySQL(gormConfig *gorm.Config) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(config.DatabaseUrl), gormConfig)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func openPostgres(gormConfig *gorm.Config) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(config.DatabaseUrl), gormConfig)
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

func migrate() error {
	return DB.AutoMigrate(&User{}, &Vault{}, &AuditLog{})
}
