package redisConfig

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"

	configErrors "github.com/mathandcrypto/cryptomath-go-captcha/internal/common/errors/config"
)

type Config struct {
	Host	string	`mapstructure:"REDIS_HOST" validate:"required"`
	Port	int16	`mapstructure:"REDIS_PORT" validate:"required,gte=1024,lte=49151"`
	Database	int	`mapstructure:"REDIS_DATABASE" validate:"gte=0,lte=15"`
}

func (c *Config) Address() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

func New() (*Config, error) {
	redisViper := viper.New()
	redisValidate := validator.New()

	redisViper.SetDefault("REDIS_HOST", "127.0.0.1")
	redisViper.SetDefault("REDIS_PORT", 6379)
	redisViper.SetDefault("REDIS_DATABASE", 1)

	redisViper.AddConfigPath("configs/redis")
	redisViper.SetConfigName("config")
	redisViper.SetConfigType("env")
	redisViper.AutomaticEnv()

	if err := redisViper.ReadInConfig(); err != nil {
		return nil, &configErrors.ReadConfigError{ConfigName: "redis", ViperErr: err}
	}

	var redisConfig Config
	if err := redisViper.Unmarshal(&redisConfig); err != nil {
		return nil, &configErrors.UnmarshalError{ConfigName: "redis", ViperErr: err}
	}

	if err := redisValidate.Struct(redisConfig); err != nil {
		return nil, &configErrors.ValidationError{ConfigName: "redis", ValidateErr: err}
	}

	return &redisConfig, nil
}