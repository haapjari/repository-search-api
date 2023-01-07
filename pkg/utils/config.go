package utils

import (
	"fmt"

	"github.com/spf13/viper"
)

// Single Point in Program to Fetch all the Environment Variables.

func GetRepositoryApiBaseUrl() string {
	viper.SetConfigFile(".env")
	viper.ReadInConfig()

	return fmt.Sprint(viper.Get("REPOSITORY_API_BASEURL"))
}

func GetLocalenv() string {
	viper.SetConfigFile(".env")
	viper.ReadInConfig()

	return fmt.Sprint(viper.Get("LOCAL_ENV"))
}

func GetDatabaseUser() interface{} {
	viper.SetConfigFile(".env")
	viper.ReadInConfig()

	return viper.Get("POSTGRES_USER")
}

func GetBaseurl() string {
	viper.SetConfigFile(".env")
	viper.ReadInConfig()

	return fmt.Sprint(viper.Get("BASEURL"))
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

func GetTempGoPath() string {
	viper.SetConfigFile("env")
	viper.ReadInConfig()

	return fmt.Sprint(viper.Get("TEMP_GOPATH"))
}

func GetGoPath() string {
	viper.SetConfigFile("env")
	viper.ReadInConfig()

	return fmt.Sprint(viper.Get("GOPATH"))
}

func GetDatabasePassword() interface{} {
	viper.SetConfigFile(".env")
	viper.ReadInConfig()

	return viper.Get("POSTGRES_PASSWORD")
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
