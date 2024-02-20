package service

import (
	"golang.org/x/crypto/bcrypt"
)

func (s *Service) ValidateUser(username, password string) (bool, error) {
	// get data from database to validate password
	user, err := s.storeClient.GetUser(username)
	if err != nil {
		return false, err
	}

	// compare stored hashed password with provided password
	err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password))
	if err != nil { // error can be ignored. Will just return false to handle response in handler
		return false, nil
	}

	return true, nil
}
