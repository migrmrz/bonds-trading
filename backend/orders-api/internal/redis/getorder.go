package redis

import (
	"strconv"
	"strings"

	"github.com/gomodule/redigo/redis"
)

func (os OrdersStore) GetOrderDetails(hashKey string) (Order, error) {
	// get redis connection
	conn := os.redis.Get()
	defer conn.Close()

	// get and return order details
	return os.hgetallFromRedis(conn, hashKey)
}

func (os OrdersStore) hgetallFromRedis(conn redis.Conn, key string) (Order, error) {
	order := Order{
		BondID: strings.Split(key, ":")[1],
	}

	arr, err := redis.Values(redis.DoWithTimeout(conn, os.redisTimeout, "hgetall", key))
	if err == redis.ErrNil {
		return order, nil
	} else if err != nil {
		return Order{}, err
	}

	keys, _ := redis.StringMap(arr, nil)

	// map keys and values to the order struct.
	order.OrderID = keys["order_id"]
	order.Action = keys["action"]
	order.Quantity, _ = strconv.Atoi(keys["quantity"])
	priceFloat64, _ := strconv.ParseFloat(keys["price"], 32)
	order.Price = float32(priceFloat64)
	order.Status = keys["status"]
	order.User = keys["user"]
	order.CreatedAt, _ = strconv.ParseInt(keys["created_at"], 10, 64)

	return order, err
}
