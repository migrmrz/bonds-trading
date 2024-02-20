package rest

// type getActiveOrdersByUserRequest struct{}

// type getActiveOrdersByUserResponse struct{}

// // StatusCode func to implement StatusCoder interface
// func (ao getActiveOrdersByUserRequest) StatusCode() int {
// 	return ao.StatusCode()
// }

// // decode request function
// func decodeGetActiveOrdersByUserRequest(_ context.Context, req *http.Request) (request interface{}, err error) {
// 	// todo get user from request

// 	getActiveOrdersByUserRequest := getActiveOrdersByUserRequest{}

// 	return getActiveOrdersByUserRequest, nil
// }

// // encode response and error functions
// func encodeGetActiveOrdersByUserResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
// 	resp := restResponse{
// 		statusCode: http.StatusOK,
// 		Data:       response,
// 	}

// 	return gokithttp.EncodeJSONResponse(ctx, w, resp)
// }

// func encodeGetActiveOrdersByUserError(ctx context.Context, err error, w http.ResponseWriter) {
// 	resp := restResponse{
// 		Data:       nil,
// 		Error:      err.Error(),
// 		statusCode: http.StatusInternalServerError,
// 	}

// 	if statusCoder, ok := err.(gokithttp.StatusCoder); ok {
// 		resp.statusCode = statusCoder.StatusCode()
// 	}

// 	if encodeErr := gokithttp.EncodeJSONResponse(ctx, w, resp); encodeErr != nil {
// 		logrus.WithFields(logrus.Fields{
// 			"func":  "errorEncoder",
// 			"step":  "EncodeJSONRepsonse",
// 			"error": encodeErr,
// 		}).Error()

// 		gokithttp.DefaultErrorEncoder(ctx, encodeErr, w)
// 	}
// }
