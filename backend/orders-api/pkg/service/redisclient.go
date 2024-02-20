package service

import (
	"trading.com/mx/orders/internal/redis"
)

type redisClient interface {
	GetOrderDetails(key string) (redis.Order, error)
	GetFilteredOrderKeys(filter, value string) ([]string, error)
	UpsertOrder(order redis.Order) (string, string, error)
}
