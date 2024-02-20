package config

import (
	"fmt"

	"github.com/nleof/goyesql"
	"github.com/spf13/viper"

	"trading.com/mx/orders/internal/redis"
	"trading.com/mx/orders/internal/store"
)

type OrdersAPIConfig struct {
	Port    int            `mapstructure:"port"`
	DBStore store.DBConfig `mapstructure:"database"`
	Redis   redis.Config   `mapstructure:"redis"`
	Queries goyesql.Queries
}

func GetConfig(configPath string) (OrdersAPIConfig, error) {
	var config OrdersAPIConfig

	queries, err := goyesql.ParseFile(configPath + "/statements.sql")
	if err != nil {
		return OrdersAPIConfig{}, fmt.Errorf("no sql file to parse")
	}

	viper.SetConfigName("orders-api")
	viper.AddConfigPath(configPath)

	err = viper.ReadInConfig()
	if err != nil {
		return config, err
	}

	err = viper.Unmarshal(&config)

	config.Queries = queries

	return config, err
}
