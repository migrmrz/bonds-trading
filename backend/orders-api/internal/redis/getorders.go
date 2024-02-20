package redis

import (
	"github.com/gomodule/redigo/redis"
)

func (os OrdersStore) GetFilteredOrderKeys(filterKey, filterValue string) ([]string, error) {
	// get redis connection.
	conn := os.redis.Get()
	defer conn.Close()

	keys, err := os.getFilteredOrderKeysFromRedis(conn, filterKey, filterValue)
	if err != nil {
		return nil, err
	}

	return keys, nil
}

func (os OrdersStore) getFilteredOrderKeysFromRedis(conn redis.Conn, filterKey, filterValue string) ([]string, error) {
	iter := 0
	keys := []string{}

	// scan for all keys until reaching the end.
	for {
		arr, err := redis.Values(redis.DoWithTimeout(
			conn, os.redisTimeout, "scan", iter, "match", "bond:*", "type", "hash", "count", 100,
		))
		if err == redis.ErrNil {
			return nil, nil
		} else if err != nil {
			return nil, err
		}

		iter, _ = redis.Int(arr[0], nil)
		k, _ := redis.Strings(arr[1], nil)

		if filterKey != "" && filterValue != "" {
			// filter out non active orders before appending to result.
			k, _ = os.filterOrdersByKey(conn, k, filterKey, filterValue)
		}

		keys = append(keys, k...)

		if iter == 0 {
			break
		}
	}

	return keys, nil
}

func (os OrdersStore) filterOrdersByKey(
	conn redis.Conn, hashKeys []string, key, filter string,
) ([]string, error) {
	var filteredKeys []string

	// get value from key and validate against filter value.
	for _, hKey := range hashKeys {
		currVal, err := redis.String(redis.DoWithTimeout(conn, os.redisTimeout, "hget", hKey, key))
		if err != nil {
			return []string{}, err
		}
		if currVal == filter {
			filteredKeys = append(filteredKeys, hKey)
		}
	}

	return filteredKeys, nil
}
