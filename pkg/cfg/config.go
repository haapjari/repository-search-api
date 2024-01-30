package utils

import (
	"github.com/spf13/viper"
)

type Config struct {
	PostgresUser     string
	PostgresPassword string
	PostgresDB       string
	PostgresHost     string
	PostgresPort     string
	GitHubUsername   string
	GitHubToken      string
}

func NewConfig() *Config {
	viper.SetConfigFile(".env")
	viper.ReadInConfig()

	return &Config{
		// viper.Get("POSTGRES_USER")
		// viper.Get("POSTGRES_PASSWORD")
		// viper.Get("POSTGRES_DB")
		// viper.Get("POSTGRES_HOST")
		// viper.Get("POSTGRES_PORT")
		// viper.Get("GITHUB_TOKEN")
		// viper.Get("GITHUB_USERNAME")

		PostgresUser:     viper.GetString("POSTGRES_USER"),
		PostgresPassword: viper.GetString(""),
		PostgresDB:       "",
		PostgresHost:     "",
		PostgresPort:     "",
		GitHubUsername:   "",
		GitHubToken:      "",
	}

}

func GetPostgresUser() interface{} {
	viper.SetConfigFile(".env")
	viper.ReadInConfig()

	return viper.Get("POSTGRES_USER")
}

func GetPostgresPassword() interface{} {
	viper.SetConfigFile(".env")
	viper.ReadInConfig()

	return viper.Get("POSTGRES_PASSWORD")
}

func GetPostgresDB() interface{} {
	viper.SetConfigFile(".env")
	viper.ReadInConfig()

	return viper.Get("POSTGRES_DB")
}

func GetDatabaseHost() interface{} {
	viper.SetConfigFile(".env")
	viper.ReadInConfig()

	return viper.Get("POSTGRES_HOST")
}

func GetDatabasePort() interface{} {
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()
	if err != nil {
		return nil
	}

	return viper.Get("POSTGRES_PORT")
}

func GetGithubToken() interface{} {
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()
	if err != nil {
		return nil
	}

	return viper.Get("GITHUB_TOKEN")
}

func GetGithubUsername() interface{} {
	viper.SetConfigFile(".env")
	viper.ReadInConfig()

	return viper.Get("GITHUB_USERNAME")
}
