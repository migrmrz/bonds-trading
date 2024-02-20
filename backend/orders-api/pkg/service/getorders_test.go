package service

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"trading.com/mx/orders/internal/redis"
)

func TestGetOrders(t *testing.T) {
	cases := []struct {
		name           string
		filter         string
		value          string
		expectedOrders []redis.Order
		expectedError  error
	}{
		{
			name:   "success",
			filter: "status",
			value:  "active",
			expectedOrders: []redis.Order{
				{
					OrderID:   "1",
					Quantity:  10,
					Action:    "sell",
					Price:     1000,
					Status:    "active",
					User:      "admin",
					CreatedAt: 1704970813,
				},
				{
					OrderID:   "2",
					Quantity:  10,
					Action:    "sell",
					Price:     1000,
					Status:    "active",
					User:      "user1",
					CreatedAt: 1704970813,
				},
			},
			expectedError: nil,
		},
	}

	redis := &redisClientMock{}
	store := &ordersClientMock{}

	srv := New(store, redis)

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			actualOrders, actualError := srv.GetOrders(tc.filter, tc.value)
			assert.Equal(t, tc.expectedOrders, actualOrders)
			assert.Equal(t, tc.expectedError, actualError)
		})
	}
}
