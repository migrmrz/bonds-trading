package service

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
	"trading.com/mx/orders/internal/redis"
	"trading.com/mx/orders/internal/store"
)

type redisClientMock struct {
	upsertOperation string
	upsertError     error
}

func (m *redisClientMock) GetOrderDetails(key string) (redis.Order, error) {
	switch key {
	case "bond:20240101000000":
		return redis.Order{
			OrderID:   "1",
			Quantity:  10,
			Action:    "sell",
			Price:     1000,
			Status:    "active",
			User:      "admin",
			CreatedAt: 1704970813,
		}, nil
	case "bond:20240102000000":
		return redis.Order{
			OrderID:   "2",
			Quantity:  10,
			Action:    "sell",
			Price:     1000,
			Status:    "active",
			User:      "user1",
			CreatedAt: 1704970813,
		}, nil
	case "bond:20240103000000":
		return redis.Order{
			OrderID:   "3",
			Quantity:  10,
			Action:    "sell",
			Price:     1000,
			Status:    "cancelled",
			User:      "admin",
			CreatedAt: 1704970813,
		}, nil
	default:
		return redis.Order{}, fmt.Errorf("key not found")
	}
}

func (m *redisClientMock) GetFilteredOrderKeys(filter, value string) ([]string, error) {
	switch filter {
	case redis.StatusFilter:
		return []string{"bond:20240101000000", "bond:20240102000000"}, nil
	case redis.UserFilter:
		return []string{}, nil
	default:
		return []string{"bond:20240101000000", "bond:20240103000000"}, nil
	}
}

func (m *redisClientMock) UpsertOrder(order redis.Order) (string, string, error) {
	switch order.BondID {
	case "bond:20240101000000": // update
		m.upsertOperation = redis.UpdateOperation
		return redis.UpdateOperation, order.OrderID, nil
	case "bond:20240104000000": // create
		m.upsertOperation = redis.InsertOperation
		return redis.InsertOperation, order.OrderID, nil
	}

	m.upsertError = fmt.Errorf("unknown error")

	return "", "", fmt.Errorf("unknown error")
}

type ordersClientMock struct{}

func (m *ordersClientMock) GetUser(username string) (store.User, error) {
	if username == "error-username" {
		return store.User{}, fmt.Errorf("an error ocurred while getting user info")
	}

	pass := []byte("admin.123")

	hash, _ := bcrypt.GenerateFromPassword(pass, bcrypt.DefaultCost)

	return store.User{
		UserID:         1,
		Username:       username,
		HashedPassword: string(hash),
		Email:          "admin@example.com",
		CreatedAt:      "1704970813",
	}, nil
}
