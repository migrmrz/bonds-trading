package service

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateUser(t *testing.T) {
	cases := []struct {
		name             string
		username         string
		password         string
		expectedResponse bool
		expectedError    error
	}{
		{
			name:             "success",
			username:         "admin",
			password:         "admin.123",
			expectedResponse: true,
			expectedError:    nil,
		},
		{
			name:             "password-not-valid",
			username:         "admin",
			password:         "manager.1",
			expectedResponse: false,
			expectedError:    nil,
		},
		{
			name:             "error-getting-user",
			username:         "error-username",
			password:         "s0methingf4ncy",
			expectedResponse: false,
			expectedError:    fmt.Errorf("an error ocurred while getting user info"),
		},
	}

	redis := &redisClientMock{}
	store := &ordersClientMock{}

	srv := New(store, redis)

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			actualResponse, actualError := srv.ValidateUser(tc.username, tc.password)
			assert.Equal(t, tc.expectedResponse, actualResponse)
			assert.Equal(t, tc.expectedError, actualError)
		})
	}
}
