package redis

import (
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/google/uuid"
)

const (
	UserFilter      = "user"
	StatusFilter    = "status"
	InsertOperation = "insert"
	UpdateOperation = "update"
)

type Order struct {
	OrderID   string  `json:"order_id"`
	BondID    string  `json:"bond_id"`
	Action    string  `json:"action"`
	Quantity  int     `json:"quantity"`
	Price     float32 `json:"price"`
	Status    string  `json:"status"`
	User      string  `json:"user"`
	CreatedAt int64   `json:"created_at"`
}

type OrdersStore struct {
	redis         *redis.Pool
	redisTimeout  time.Duration
	nowFunc       func() time.Time
	uuidFunc      func() string
	keyExpiration int
}

type Config struct {
	Address       string        `mapstructure:"address"`
	MaxActive     int           `mapstructure:"max-active"`
	MaxIdle       int           `mapstructure:"max-idle"`
	IdleTimeout   time.Duration `mapstructure:"idle-timeout"`
	MaxLifetime   time.Duration `mapstructure:"max-lifetime"`
	CallTimeout   time.Duration `mapstructure:"call-timeout"`
	KeyExpiration int           `mapstructure:"key-expiration"`
}

func NewStorage(config Config) (OrdersStore, error) {
	pool := &redis.Pool{
		MaxIdle:     config.MaxIdle,
		MaxActive:   config.MaxActive,
		IdleTimeout: config.IdleTimeout,
		Wait:        true,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", config.Address)
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}

	return OrdersStore{
		redis:        pool,
		redisTimeout: config.CallTimeout,
		nowFunc:      time.Now,
		uuidFunc: func() string {
			return uuid.New().String()
		},
		keyExpiration: config.KeyExpiration,
	}, nil
}
