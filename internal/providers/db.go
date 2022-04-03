package providers

import (
	"database/sql"
	"fmt"

	_ "github.com/volatiletech/sqlboiler/v4/drivers/sqlboiler-psql/driver"

	dbConfig "github.com/mathandcrypto/cryptomath-go-captcha/configs/db"
)

func NewDBProvider(config *dbConfig.Config) (*sql.DB, error) {
	db, err := sql.Open("postgres", config.DSN())

	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database connection: %w", err)
	}

	return db, nil
}