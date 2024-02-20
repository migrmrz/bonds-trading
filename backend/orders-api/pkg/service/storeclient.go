package service

import "trading.com/mx/orders/internal/store"

type storeClient interface {
	GetUser(username string) (store.User, error)
}
