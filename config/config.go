package config

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Dsn                 string    `mapstructure:"DSN"`
	Port                string    `mapstructure:"PORT"`
	TokenSymmetricKey   string    `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	AccessTokenDuration time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	JwtSecretKey        string    `mapstructure:"JWT_SECRET_KEY"`
	RefreshJwtSecretKey string    `mapstructure:"REFRESH_JWT_SECRET_KEY"`
	Email               string    `mapstructure:"EMAIL"`
	EmailPassword       string    `mapstructure:"EMAIL_PASSWORD"`
}

func LoadConfig(path string) (config *Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return config, err
	}

	err = viper.Unmarshal(&config)
	return config, err
}
