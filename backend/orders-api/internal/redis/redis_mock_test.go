package redis

import (
	"fmt"
	"time"
)

type RedisConnMock struct {
	calls [][]interface{}
}

func (rc *RedisConnMock) DoWithTimeout(_ time.Duration, command string, args ...interface{}) (interface{}, error) {
	rc.calls = append(rc.calls, append([]interface{}{command}, args...))

	arg0 := ""
	if len(args) > 0 {
		arg0 = args[0].(string)
	}

	switch command {
	case "hset":
		switch arg0 {
		case "bond:20240101000000":
			return nil, nil
		case "bond:20240102000000":
			return nil, nil
		default:
			return nil, fmt.Errorf("key not found")
		}
	case "expire":
		switch arg0 {
		case "bond:20240101000000":
			return nil, nil
		case "bond:20240102000000":
			return nil, nil
		default:
			return nil, fmt.Errorf("key not found")
		}
	case "hgetall":
		switch arg0 {
		case "bond:20240107152312":
			return `{"action":"sell","quantity":3,"price":2000,"status":"active","user:"admin","created_at":1704766409}`,
				nil
		case "bond:20240106122105":
			return nil, nil
		default:
			return nil, nil
		}

	case "exists":
		switch arg0 {
		case "bond:20240101000000":
			return int64(0), nil
		default:
			return int64(1), nil
		}

	case "hget":
		switch arg0 {
		case "bond:20240102000000":
			return "a87a260f-50d1-41c1-a081-9cc58c1cb075", nil
		default:
			return "", fmt.Errorf("error getting values")
		}

	default:
		return nil, fmt.Errorf("unexpected redis command")
	}
}

func (rc RedisConnMock) Close() error {
	return nil
}

func (rc RedisConnMock) Err() error {
	return nil
}

func (rc RedisConnMock) Do(string, ...interface{}) (interface{}, error) {
	return nil, nil
}

func (rc RedisConnMock) Send(string, ...interface{}) error {
	return nil
}

func (rc RedisConnMock) Flush() error {
	return nil
}

func (rc RedisConnMock) Receive() (interface{}, error) {
	return nil, nil
}

func (rc RedisConnMock) ReceiveWithTimeout(time.Duration) (interface{}, error) {
	return nil, nil
}
