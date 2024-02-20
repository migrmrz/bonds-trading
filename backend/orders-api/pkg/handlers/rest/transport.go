package rest

import (
	"context"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-kit/kit/endpoint"
	gokithttp "github.com/go-kit/kit/transport/http"
	"github.com/sirupsen/logrus"

	"trading.com/mx/orders/internal/redis"
)

// Defining endpoint for the service
type Endpoints struct {
	GetOrders   endpoint.Endpoint
	UpsertOrder endpoint.Endpoint
}

type OrdersHandler interface {
	ValidateUser(username, password string) (bool, error)
	GetOrders(filter, value string) ([]redis.Order, error)
	UpsertOrder(order redis.Order)
}

type userData struct {
	username string
	password string
}

type restResponse struct {
	statusCode int
	Error      string      `json:"error,omitempty"`
	Data       interface{} `json:"data"`
}

func (rr restResponse) StatusCode() int {
	return rr.statusCode
}

func MakeHTTPHandlers(s OrdersHandler) http.Handler {
	r := chi.NewRouter()

	version1RouterFunc := makeVersion1RouterFunc(s)
	r.Route("/book/v1", version1RouterFunc)

	return r
}

func makeVersion1RouterFunc(s OrdersHandler) func(chi.Router) {
	endpoints := makeServerEndpoints(s)

	insertOrderServer := gokithttp.NewServer(
		endpoints.UpsertOrder,
		decodeInsertOrderRequest,
		encodeResponse,
		gokithttp.ServerBefore(gokithttp.PopulateRequestContext),
		gokithttp.ServerErrorEncoder(encodeError),
	)

	getActiveOrdersServer := gokithttp.NewServer(
		endpoints.GetOrders,
		decodeGetOrdersRequest,
		encodeResponse,
		gokithttp.ServerBefore(gokithttp.PopulateRequestContext),
		gokithttp.ServerErrorEncoder(encodeError),
	)

	return func(r chi.Router) {
		r.Put("/orders", insertOrderServer.ServeHTTP)
		r.Get("/orders", getActiveOrdersServer.ServeHTTP)
	}
}

func makeServerEndpoints(s OrdersHandler) Endpoints {
	return Endpoints{
		UpsertOrder: makeInsertOrderEndpoint(s),
		GetOrders:   makeGetOrdersEndpoint(s),
	}
}

// generic response
func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	resp := restResponse{
		Data:       response,
		statusCode: http.StatusOK,
	}

	return gokithttp.EncodeJSONResponse(ctx, w, resp)
}

func encodeError(ctx context.Context, err error, w http.ResponseWriter) {
	resp := restResponse{
		Data:       nil,
		Error:      err.Error(),
		statusCode: http.StatusInternalServerError, // default status code
	}

	if statusCoder, ok := err.(gokithttp.StatusCoder); ok {
		resp.statusCode = statusCoder.StatusCode()
	}

	if encodeErr := gokithttp.EncodeJSONResponse(ctx, w, resp); encodeErr != nil {
		logrus.WithFields(logrus.Fields{
			"function": "errorEncoder",
			"step":     "EncodeJSONResponse",
			"error":    encodeErr,
		}).Error()
		gokithttp.DefaultErrorEncoder(ctx, encodeErr, w)
	}
}
