package databaseConfig

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"

	configErrors "github.com/mathandcrypto/cryptomath-go-captcha/internal/common/errors/config"
)

type Config struct {
	Host	string	`mapstructure:"DATABASE_HOST" validate:"required"`
	Port	int16	`mapstructure:"DATABASE_PORT" validate:"required"`
	User	string	`mapstructure:"POSTGRES_USER" validate:"required"`
	Password	string	`mapstructure:"POSTGRES_PASSWORD" validate:"required"`
	Database	string	`mapstructure:"POSTGRES_DB" validate:"required"`
}

func (c *Config) DSN() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable", c.Host, c.User, c.Password, c.Database, c.Port)
}

func New() (*Config, error) {
	dbViper := viper.New()
	dbValidate := validator.New()

	dbViper.SetDefault("DATABASE_HOST", "localhost")
	dbViper.SetDefault("DATABASE_PORT", 5432)

	dbViper.AddConfigPath("configs/database")
	dbViper.SetConfigName("config")
	dbViper.SetConfigType("env")
	dbViper.AutomaticEnv()

	if err := dbViper.ReadInConfig(); err != nil {
		return nil, &configErrors.ReadConfigError{ConfigName: "database", ViperErr: err}
	}

	var dbConfig Config
	if err := dbViper.Unmarshal(&dbConfig); err != nil {
		return nil, &configErrors.UnmarshalError{ConfigName: "database", ViperErr: err}
	}

	if err := dbValidate.Struct(dbConfig); err != nil {
		return nil, &configErrors.ValidationError{ConfigName: "database", ValidateErr: err}
	}

	return &dbConfig, nil
}
