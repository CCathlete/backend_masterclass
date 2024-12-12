package sqlc

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/jackc/pgx/v4/stdlib"
)

type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// Make sure to db.Close()!
func ConnectToDB() (*sql.DB, error) {
	db, err := sql.Open("pgx", connectionString())
	if err != nil {
		return nil, fmt.Errorf("error when opening db: %w", err)
	}

	return db, nil
}

func connectionString() string {
	cfg := PostgresConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASS"),
		DBName:   os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("DB_SSL_MODE"),
	}

	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User,
		cfg.Password, cfg.DBName, cfg.SSLMode,
	)
}
