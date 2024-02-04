package cfg

import (
	"github.com/spf13/viper"
)

type Config struct {
	Port    string
	GinMode string
}

func NewConfig() *Config {
	viper.SetConfigFile(".env")
	viper.ReadInConfig()

	return &Config{
		Port:    viper.GetString("PORT"),
		GinMode: viper.GetString("GIN_MODE"),
	}
}
