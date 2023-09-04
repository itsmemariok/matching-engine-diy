package main

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

func writeIndexNamesConfig(indexName, createdIndexEndpointName string) {
	viper.SetConfigName("config.toml")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")

	viper.SetDefault("indexName", indexName)

	parts := strings.Split(indexName, "/")
	indexID := parts[len(parts)-1]

	viper.SetDefault("indexid", indexID)

	parts = strings.Split(createdIndexEndpointName, "/")
	indexEndpointID := parts[len(parts)-1]

	viper.SetDefault("indexendpointid", indexEndpointID)

	err := viper.WriteConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
}

func readIndexNamesConfig() (string, string) {
	viper.SetConfigName("config.toml")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	return viper.GetString("indexName"), viper.GetString("indexendpointid")
}

func writeDeployedIndexEndpointURL(url string) {
	viper.SetConfigName("config.toml")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")

	viper.SetDefault("deployedIndexEndpointURL", url)

	err := viper.WriteConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
}

func readDeployedIndexEndpointURL() string {
	viper.SetConfigName("config.toml")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	return viper.GetString("deployedIndexEndpointURL")
}

func readConfig() GlobalConfig {
	viper.SetConfigName("config.toml")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	var cfg GlobalConfig

	err = viper.Unmarshal(&cfg)
	if err != nil {
		panic(fmt.Errorf("unmarshal error: %s \n", err))
	}

	return cfg
}
