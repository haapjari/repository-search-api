package cfg

import (
	"github.com/spf13/viper"
	"log/slog"
)

type Config struct {
	Port        string
	EnablePprof bool
}

const (
	PortKey        = "PORT"
	EnablePprofKey = "ENABLE_PPROF"
)

func NewConfig() *Config {
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()
	if err != nil {
		slog.Warn("unable to read the configuration file: " + err.Error())
	}

	viper.AutomaticEnv()

	return &Config{
		Port:        viper.GetString(PortKey),
		EnablePprof: viper.GetBool(EnablePprofKey),
	}
}
