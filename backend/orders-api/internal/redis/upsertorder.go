package redis

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/gomodule/redigo/redis"
)

func (os OrdersStore) UpsertOrder(order Order) (opRes, orderID string, err error) {
	// get redis connection.
	conn := os.redis.Get()
	defer conn.Close()

	key := fmt.Sprintf("bond:%s", order.BondID)

	operation, orderID, err := os.hsetValuesInRedis(conn, key, order)
	if err != nil {
		return "", "", err
	}

	return operation, orderID, nil
}

func (os OrdersStore) hsetValuesInRedis(conn redis.Conn, key string, order Order) (opRes, orderID string, err error) {
	// check if key already exists.
	res, err := redis.DoWithTimeout(conn, os.redisTimeout, "exists", key)
	if err != nil {
		return "", "", err
	}

	exists, _ := redis.Int64(res, nil)
	if exists == 1 {
		// update order: get existing order number and set it in Order struct.
		res, err = redis.DoWithTimeout(conn, os.redisTimeout, "hget", key, "order_id")
		if err != nil {
			return "", "", err
		}

		orderID, _ = redis.String(res, nil)

		order.OrderID = orderID

		opRes = UpdateOperation
	} else {
		// create order: generate uuid and set it in Order struct.
		order.OrderID = os.uuidFunc()
		opRes = InsertOperation
	}

	v := reflect.ValueOf(order)

	values := make([]interface{}, 0, v.NumField()*2)

	// convert keys and values from Order{} into a list for further processsing.
	for i := 0; i < v.NumField(); i++ {
		if v.Type().Field(i).Tag == `json:"bond_id"` || // already part of the key
			// opRes == UpdateOperation && v.Type().Field(i).Tag == `json:"order_id"` ||
			opRes == UpdateOperation && v.Type().Field(i).Tag == `json:"created_at"` { // do not update
			continue
		}

		// get json tag.
		values = append(values, strings.ReplaceAll(strings.Split(string(v.Type().Field(i).Tag), ":")[1], `"`, ""))
		// get value
		values = append(values, v.Field(i).Interface())
	}

	// insert key with fields
	_, err = redis.DoWithTimeout(
		conn,
		os.redisTimeout,
		"hset",
		append([]interface{}{key}, values...)...,
	)
	if err != nil {
		return "", orderID, err
	}

	// set expiration for key
	_, err = redis.DoWithTimeout(conn, os.redisTimeout, "expire", key, os.keyExpiration)
	if err != nil {
		return "", orderID, err
	}

	return opRes, orderID, nil

}
