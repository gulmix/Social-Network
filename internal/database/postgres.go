package database

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/gulmix/Social-Network/internal/config"
)

var DB *sql.DB

func InitPostgres(cfg *config.Config) error {
	dsn := cfg.Database.ConnectionString()

	var err error
	DB, err = sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	if err = DB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Successfully connected to PostgreSQL")
	return nil
}

func ClosePostgres() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}
