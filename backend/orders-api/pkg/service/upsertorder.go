package service

import (
	"encoding/json"

	nats "github.com/nats-io/nats.go"
	"github.com/sirupsen/logrus"
	"trading.com/mx/orders/internal/redis"
)

func (s *Service) UpsertOrder(order redis.Order) {
	operation, orderID, err := s.redisClient.UpsertOrder(order)
	if err != nil {
		logrus.Error(logrus.WithFields(logrus.Fields{
			"error":        err.Error(),
			"operation":    operation,
			"orderDetails": order,
		}))
	}

	order.OrderID = orderID

	if operation == redis.UpdateOperation {
		// create local connection to NATS.
		nc, err := nats.Connect(nats.DefaultURL)
		if err != nil {
			logrus.Error(logrus.WithFields(logrus.Fields{
				"error":        err.Error(),
				"operation":    operation,
				"orderDetails": order,
			}))
		}

		// close connection on function exit.
		defer nc.Close()

		// Marshal JSON for publishing
		orderJson, err := json.Marshal(order)
		if err != nil {
			logrus.Error(logrus.WithFields(logrus.Fields{
				"error":        err.Error(),
				"orderDetails": order,
			}))
		}

		// Publish message.
		err = nc.Publish("com.trading.orders.update", orderJson)
		if err != nil {
			logrus.Error(logrus.WithFields(logrus.Fields{
				"error":        err.Error(),
				"orderDetails": order,
			}))
		}
	}
}
