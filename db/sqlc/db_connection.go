package sqlc

import (
	u "backend-masterclass/util"
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v4/stdlib"
)

// Make sure to db.Close()!
func ConnectToDB(cfg u.Config) (*sql.DB, error) {

	db, err := sql.Open("pgx", connectionString(cfg))
	if err != nil {
		return nil, fmt.Errorf("error when opening db: %w", err)
	}

	return db, nil
}

func connectionString(cfg u.Config) string {

	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User,
		cfg.Password, cfg.DBName, cfg.SSLMode,
	)
}
