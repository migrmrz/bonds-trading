package config

import (
	"testing"
	"time"

	"github.com/nleof/goyesql"
	"github.com/stretchr/testify/assert"

	"trading.com/mx/orders/internal/redis"
	"trading.com/mx/orders/internal/store"
)

func TestGetConfig(t *testing.T) {
	expectedConfig := OrdersAPIConfig{
		Port: 8001,
		DBStore: store.DBConfig{
			Host:            "localhost",
			Port:            "5432",
			DBName:          "ordersdb",
			User:            "postgres",
			MaxOpenConn:     30,
			MaxIdleConn:     15,
			ConnMaxLifetime: time.Minute * 5,
		},
		Redis: redis.Config{
			Address:       "localhost:6379",
			MaxActive:     50,
			MaxIdle:       10,
			IdleTimeout:   10 * time.Minute,
			MaxLifetime:   0,
			CallTimeout:   10 * time.Second,
			KeyExpiration: 3600,
		},
		Queries: goyesql.Queries{
			"getUser": "SELECT * FROM USERS WHERE username = $1;",
		},
	}

	actualConfig, err := GetConfig(".")
	assert.Equal(t, expectedConfig, actualConfig)
	assert.NoError(t, err)
}
