package service

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"trading.com/mx/orders/internal/redis"
)

func TestUpsertOrders(t *testing.T) {
	cases := []struct {
		name             string
		order            redis.Order
		expectedResponse string
		expectedError    error
	}{
		{
			name: "success-update",
			order: redis.Order{
				OrderID:   "1",
				BondID:    "bond:20240101000000",
				Quantity:  10,
				Action:    "sell",
				Price:     1000,
				Status:    "cancelled",
				User:      "admin",
				CreatedAt: 1704970813,
			},
			expectedResponse: redis.UpdateOperation,
			expectedError:    nil,
		},
		{
			name: "success-insert",
			order: redis.Order{
				OrderID:   "1",
				BondID:    "bond:20240104000000",
				Quantity:  10,
				Action:    "sell",
				Price:     1000,
				Status:    "cancelled",
				User:      "admin",
				CreatedAt: 1704970813,
			},
			expectedResponse: redis.InsertOperation,
			expectedError:    nil,
		},
		{
			name: "error",
			order: redis.Order{
				OrderID:   "1",
				BondID:    "bond:error",
				Quantity:  10,
				Action:    "sell",
				Price:     1000,
				Status:    "active",
				User:      "admin",
				CreatedAt: 1704970813,
			},
			expectedResponse: "",
			expectedError:    fmt.Errorf("unknown error"),
		},
	}

	redis := &redisClientMock{}
	store := &ordersClientMock{}

	srv := New(store, redis)

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			redis.upsertError = nil
			redis.upsertOperation = ""

			srv.UpsertOrder(tc.order)
			assert.Equal(t, tc.expectedResponse, redis.upsertOperation)
			assert.Equal(t, tc.expectedError, redis.upsertError)
		})
	}
}
