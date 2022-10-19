package cfg

import (
	"github.com/spf13/viper"
)

type Config struct {
	Port          string
	QueryInterval string
}

func NewConfig() (*Config, error) {
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	return &Config{
		Port:          viper.GetString("PORT"),
		QueryInterval: viper.GetString("QUERY_INTERVAL"),
	}, nil
}
