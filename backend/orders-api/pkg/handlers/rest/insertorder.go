package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/google/uuid"
	"trading.com/mx/orders/internal/redis"
)

type upsertOrderRequest struct {
	userData userData
	OrderID  uuid.UUID
	BondID   string  `json:"bond_id"`
	Quantity int     `json:"quantity"`
	Action   string  `json:"action"`
	Price    float32 `json:"price"`
	Status   string  `json:"status"`
	Username string  `json:"user"`
}

func makeInsertOrderEndpoint(s OrdersHandler) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(upsertOrderRequest)

		// Validate user and password
		ok, err := s.ValidateUser(req.userData.username, req.userData.password)
		if err != nil {
			return nil, ErrInternalServerError{err.Error()}
		}
		if !ok {
			return nil, ErrUnauthorized{"wrong user and/or password provided"}
		}

		order := redis.Order{
			OrderID:   req.OrderID.String(),
			BondID:    req.BondID,
			Quantity:  req.Quantity,
			Action:    req.Action,
			Price:     req.Price,
			Status:    req.Status,
			User:      req.Username,
			CreatedAt: time.Now().Unix(),
		}

		// Insert data
		go s.UpsertOrder(order)

		return fmt.Sprintf("bond:%s", req.BondID), nil
	}
}

// decode request function
func decodeInsertOrderRequest(_ context.Context, req *http.Request) (request interface{}, err error) {
	// read body
	reqBody, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, ErrBadRequest{"can't read body"}
	}

	var upsertReq upsertOrderRequest

	// parse body
	err = json.Unmarshal(reqBody, &upsertReq)
	if err != nil {
		return nil, ErrBadRequest{"unable to parse body"}
	}

	// get data from header
	usr, pwd, ok := req.BasicAuth()
	if !ok {
		return nil, ErrUnauthorized{"no authentication data provided"}
	}

	upsertReq.userData.username = usr
	upsertReq.userData.password = pwd

	return upsertReq, nil
}
