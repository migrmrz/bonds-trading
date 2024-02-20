package redis

import (
	"testing"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/stretchr/testify/assert"
)

func TestUpsertOrder(t *testing.T) {
	cases := []struct {
		name               string
		order              Order
		expectedError      error
		expectedRedisCalls [][]interface{}
	}{
		{
			name: "insert-success",
			order: Order{
				OrderID:   "a87a260f-50d1-41c1-a081-9cc58c1cb074",
				BondID:    "20240101000000",
				Action:    "sell",
				Quantity:  10,
				Price:     100,
				Status:    "active",
				User:      "admin",
				CreatedAt: 1704942731,
			},
			expectedError: nil,
			expectedRedisCalls: [][]interface{}{
				{"exists", "bond:20240101000000"},
				{"hset", "bond:20240101000000", "order_id", "a87a260f-50d1-41c1-a081-9cc58c1cb074", "action", "sell", "quantity", 10, "price", float32(100), "status", "active", "user", "admin", "created_at", int64(1704942731)},
				{"expire", "bond:20240101000000", 3600},
			},
		},
		{
			name: "update-success",
			order: Order{
				BondID:    "20240102000000",
				Action:    "sell",
				Quantity:  10,
				Price:     100,
				Status:    "active",
				User:      "admin",
				CreatedAt: 1704942731,
			},
			expectedError: nil,
			expectedRedisCalls: [][]interface{}{
				{"exists", "bond:20240102000000"},
				{"hget", "bond:20240102000000", "order_id"},
				{"hset", "bond:20240102000000", "order_id", "a87a260f-50d1-41c1-a081-9cc58c1cb075", "action", "sell", "quantity", 10, "price", float32(100), "status", "active", "user", "admin", "created_at", int64(1704942731)},
				{"expire", "bond:20240102000000", 3600},
			},
		},
	}

	redisConnMock := &RedisConnMock{}
	store := OrdersStore{
		redis:         &redis.Pool{Dial: func() (redis.Conn, error) { return redisConnMock, nil }},
		nowFunc:       dummyNowFunc,
		uuidFunc:      dummyUUIDFunc,
		keyExpiration: 3600,
	}

	for _, c := range cases {
		redisConnMock.calls = nil
		t.Run(c.name, func(t *testing.T) {
			_, _, actualError := store.UpsertOrder(c.order)
			assert.Equal(t, c.expectedError, actualError)
			assert.Equal(t, c.expectedRedisCalls, redisConnMock.calls)
		})
	}
}

func dummyNowFunc() time.Time {
	return time.Date(2024, 01, 01, 0, 0, 1, 0, time.UTC)
}

func dummyUUIDFunc() string {
	return "a87a260f-50d1-41c1-a081-9cc58c1cb074"
}
