package service

import (
	"github.com/sirupsen/logrus"

	"trading.com/mx/orders/internal/redis"
)

func (s *Service) GetOrders(filter, value string) ([]redis.Order, error) {
	orderKeys, err := s.redisClient.GetFilteredOrderKeys(filter, value)
	if err != nil {
		return nil, err
	}

	orders := []redis.Order{}

	for _, key := range orderKeys {
		order, err := s.redisClient.GetOrderDetails(key)
		if err != nil {
			logrus.Warnf("unable to get order details: %s", err.Error())
		}

		orders = append(orders, order)

	}

	return orders, nil
}
