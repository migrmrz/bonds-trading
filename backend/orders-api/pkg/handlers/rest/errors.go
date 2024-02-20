package rest

import "net/http"

// ErrBadRequest represents a bad request error.
type ErrBadRequest struct {
	message string
}

// Error implements the error interface.
func (e ErrBadRequest) Error() string {
	return e.message
}

// StatusCode implements the StatusCoder interaface.
func (e ErrBadRequest) StatusCode() int {
	return http.StatusBadRequest
}

// ErrUnauthorized represents an unauthorized error.
type ErrUnauthorized struct {
	message string
}

// Error implements the error interface.
func (e ErrUnauthorized) Error() string {
	return e.message
}

// StatusCode implements the StatusCoder interface.
func (e ErrUnauthorized) StatusCode() int {
	return http.StatusUnauthorized
}

// ErrInternalServerError represents an error.
type ErrInternalServerError struct {
	message string
}

// Error implements the error interface.
func (e ErrInternalServerError) Error() string {
	return e.message
}

// StatusCode implements the StatusCoder interface.
func (e ErrInternalServerError) StatusCode() int {
	return http.StatusInternalServerError
}
