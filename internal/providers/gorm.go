package providers

import (
	"context"
	"fmt"
	"time"

	logger "github.com/mathandcrypto/cryptomath-gorm-logger"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	databaseConfig "github.com/mathandcrypto/cryptomath-go-captcha/configs/database"
)

func NewGORMProvider(ctx context.Context, l *logrus.Logger, config *databaseConfig.Config) (*gorm.DB, error) {
	newLogger := logger.New(l, logger.Config{
		SlowThreshold:         time.Second,
		SourceField:           "source",
		SkipErrRecordNotFound: true,
	})

	db, err := gorm.Open(postgres.Open(config.DSN()), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}

	return db.WithContext(ctx), nil
}
