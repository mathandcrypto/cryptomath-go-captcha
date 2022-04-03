package appConfig

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"

	configErrors "github.com/mathandcrypto/cryptomath-go-captcha/internal/common/errors/config"
)

type Config struct {
	Host	string	`mapstructure:"APP_HOST" validate:"required"`
	Port	int16	`mapstructure:"APP_PORT" validate:"required,gte=1024,lte=49151"`
}

func (c *Config) Address() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

func New() (*Config, error) {
	appViper := viper.New()
	appValidate := validator.New()

	appViper.SetDefault("APP_HOST", "127.0.0.1")

	appViper.AddConfigPath("configs/app")
	appViper.SetConfigName("config")
	appViper.SetConfigType("env")
	appViper.AutomaticEnv()

	if err := appViper.ReadInConfig(); err != nil {
		return nil, &configErrors.ReadConfigError{ConfigName: "app", ViperErr: err}
	}

	var appConfig Config
	if err := appViper.Unmarshal(&appConfig); err != nil {
		return nil, &configErrors.UnmarshalError{ConfigName: "app", ViperErr: err}
	}

	if err := appValidate.Struct(appConfig); err != nil {
		return nil, &configErrors.ValidationError{ConfigName: "app", ValidateErr: err}
	}

	return &appConfig, nil
}