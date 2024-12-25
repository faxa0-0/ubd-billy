package postgres

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type PostgresStorage struct {
	db *sql.DB
}

func NewPostgresStorage(dsn string) (*PostgresStorage, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("could not connect to database: %v", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}
	log.Println("DB: Successfully connected!")

	if err := createTables(db); err != nil {
		return nil, fmt.Errorf("could not create tables: %v", err)
	}

	return &PostgresStorage{db: db}, nil
}
func createTables(db *sql.DB) error {
	usersTableQuery := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		name VARCHAR(100) NOT NULL,
		login VARCHAR(100) UNIQUE NOT NULL,
		pass VARCHAR(255) NOT NULL,
		plan VARCHAR(255) NOT NULL,
		role VARCHAR(50) NOT NULL,
		active BOOLEAN DEFAULT TRUE
	);`

	if _, err := db.Exec(usersTableQuery); err != nil {
		return fmt.Errorf("could not create users table: %w", err)
	}

	usageTableQuery := `
	CREATE TABLE IF NOT EXISTS usage (
		id SERIAL PRIMARY KEY,
		user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		youtube NUMERIC DEFAULT 0,
		netflix NUMERIC DEFAULT 0,
		spotify NUMERIC DEFAULT 0,
		basic NUMERIC DEFAULT 0,
		verified_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	);`
	if _, err := db.Exec(usageTableQuery); err != nil {
		return fmt.Errorf("could not create usage table: %w", err)
	}

	return nil
}
