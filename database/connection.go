package database

import (
	"fmt"
	"log/slog"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

// Connect connects to PostgreSQL database using DSN
func Connect(host, port, user, password, dbname, sslmode string) (*sqlx.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode)

	db, err := sqlx.Connect("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Log current database name for debugging
	var currentDB string
	if err := db.Get(&currentDB, "SELECT current_database()"); err != nil {
		slog.Warn("Could not get current database name", "error", err)
	} else {
		slog.Info("Successfully connected to database", "database", currentDB)
	}

	return db, nil
}
