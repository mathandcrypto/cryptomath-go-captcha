package main

import (
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	dbConfig "github.com/mathandcrypto/cryptomath-go-captcha/configs/db"
	"github.com/mathandcrypto/cryptomath-go-captcha/internal/common/logger"
)

func main() {
	// Init logger
	l := logger.CreateLogger("captcha-migrate")

	//	Init database
	dbConf, err := dbConfig.New()
	if err != nil {
		l.WithError(err).Fatal("failed to load database config")
	}

	//	Init migration
	m, err := migrate.New("file://migrations", dbConf.URL())
	if err != nil {
		l.WithError(err).Fatal("failed to create new migrate instance")
	}

	if err = m.Up(); err != nil && err != migrate.ErrNoChange {
		l.WithError(err).Fatal("failed to up migrations")
	}
}