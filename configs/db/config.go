package db

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"

	configErrors "github.com/mathandcrypto/cryptomath-go-captcha/internal/common/errors/config"
)

type Config struct {
	Host    string   `mapstructure:"DB_HOST" validate:"required"`
	Port	int16	`mapstructure:"DB_PORT" validate:"required,gte=1024,lte=49151"`
	User	string	`mapstructure:"DB_USER" validate:"required"`
	Password	string	`mapstructure:"DB_PASSWORD" validate:"required"`
	Database	string	`mapstructure:"DB_NAME" validate:"required"`
}

func (c *Config) DSN() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable", c.Host, c.User, c.Password, c.Database, c.Port)
}

func (c *Config) URL() string {
	return fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=disable", c.User, c.Password, c.Host, c.Port, c.Database)
}

func New() (*Config, error) {
	dbViper := viper.New()
	dbValidate := validator.New()

	dbViper.SetDefault("DB_HOST", "127.0.0.1")
	dbViper.SetDefault("DB_PORT", 5432)

	dbViper.AddConfigPath("configs/db")
	dbViper.SetConfigName("config")
	dbViper.SetConfigType("env")
	dbViper.AutomaticEnv()

	if err := dbViper.ReadInConfig(); err != nil {
		return nil, &configErrors.ReadConfigError{ConfigName: "db", ViperErr: err}
	}

	var dbConfig Config
	if err := dbViper.Unmarshal(&dbConfig); err != nil {
		return nil, &configErrors.UnmarshalError{ConfigName: "db", ViperErr: err}
	}

	if err := dbValidate.Struct(dbConfig); err != nil {
		return nil, &configErrors.ValidationError{ConfigName: "db", ValidateErr: err}
	}

	return &dbConfig, nil
}