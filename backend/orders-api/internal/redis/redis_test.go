package redis

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewStorage(t *testing.T) {
	config := Config{
		Address:       "localhost",
		MaxActive:     10,
		MaxIdle:       10,
		IdleTimeout:   10,
		MaxLifetime:   10,
		CallTimeout:   10,
		KeyExpiration: 10,
	}

	emptyStore := OrdersStore{}

	actualStore, actualError := NewStorage(config)
	if reflect.DeepEqual(emptyStore, actualStore) {
		t.Error("error creating redis instance")
	}

	assert.Nil(t, actualError)
}
