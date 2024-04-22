package cfg

import (
	"github.com/spf13/viper"
)

type Config struct {
	Port          string
	QueryInterval string
}

const (
	PortKey          = "PORT"
	QueryIntervalKey = "QUERY_INTERVAL"
)

func NewConfig() (*Config, error) {
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	return &Config{
		Port:          viper.GetString(PortKey),
		QueryInterval: viper.GetString(QueryIntervalKey),
	}, nil
}
