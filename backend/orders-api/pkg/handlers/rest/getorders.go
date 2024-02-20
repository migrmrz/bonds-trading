package rest

import (
	"context"
	"net/http"

	"github.com/go-kit/kit/endpoint"
	"trading.com/mx/orders/internal/redis"
)

type filter struct {
	key   string
	value string
}

type getOrdersRequest struct {
	UserData userData
	filter   filter
}

func makeGetOrdersEndpoint(s OrdersHandler) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(getOrdersRequest)

		// Validate user and password
		ok, err := s.ValidateUser(req.UserData.username, req.UserData.password)
		if err != nil {
			return nil, ErrInternalServerError{err.Error()}
		}
		if !ok {
			return nil, ErrUnauthorized{"wrong user and/or password provided"}
		}

		orders, err := s.GetOrders(req.filter.key, req.filter.value)
		if err != nil {
			return nil, err
		}

		return orders, nil
	}
}

// decode request function
func decodeGetOrdersRequest(_ context.Context, req *http.Request) (request interface{}, err error) {
	var getReq getOrdersRequest

	// get user data from header.
	usr, pwd, ok := req.BasicAuth()
	if !ok {
		return nil, ErrUnauthorized{"no authentication data provided"}
	}

	getReq.UserData.username = usr
	getReq.UserData.password = pwd

	// get filter (if any) from query params for the order search.
	params := req.URL.Query()

	switch {
	case params.Has(redis.UserFilter):
		getReq.filter.key = redis.UserFilter
		getReq.filter.value = params.Get(redis.UserFilter)
	case params.Has(redis.StatusFilter):
		getReq.filter.key = redis.StatusFilter
		getReq.filter.value = params.Get(redis.StatusFilter)
	default:
		getReq.filter.key, getReq.filter.value = "", ""
	}

	return getReq, nil
}
