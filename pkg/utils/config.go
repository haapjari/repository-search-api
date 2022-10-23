package utils

import (
	"fmt"

	"github.com/spf13/viper"
)

// Single Point in Program to Fetch all the Environment Variables.

func GetDatabaseUser() interface{} {
	viper.SetConfigFile(".env")
	viper.ReadInConfig()

	return viper.Get("POSTGRES_USER")
}

func GetGithubGraphQlApiBaseurl() string {
	viper.SetConfigFile("env")
	viper.ReadInConfig()

	return fmt.Sprint(viper.Get("GITHUB_GRAPHQL_API_BASEURL"))
}

func GetSourceGraphGraphQlApiBaseurl() string {
	viper.SetConfigFile("env")
	viper.ReadInConfig()

	return fmt.Sprint(viper.Get("SOURCEGRAPH_GRAPHQL_API_BASEURL"))
}

func GetDatabasePassword() interface{} {
	viper.SetConfigFile(".env")
	viper.ReadInConfig()

	return viper.Get("POSTGRES_PASSWORD")
}

func GetPrimaryLanguage() string {
	viper.SetConfigFile(".env")
	viper.ReadInConfig()

	return fmt.Sprint(viper.Get("PRIMARY_LANGUAGE"))
}

func GetKeyword() string {
	viper.SetConfigFile(".env")
	viper.ReadInConfig()

	return fmt.Sprint(viper.Get("KEYWORD"))
}

func GetDatabaseName() interface{} {
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
	viper.ReadInConfig()

	return viper.Get("POSTGRES_PORT")
}

func GetGithubApiToken() interface{} {
	viper.SetConfigFile(".env")
	viper.ReadInConfig()

	return viper.Get("GITHUB_API_TOKEN")
}

func GetGithubUsername() interface{} {
	viper.SetConfigFile(".env")
	viper.ReadInConfig()

	return viper.Get("GITHUB_USERNAME")
}
