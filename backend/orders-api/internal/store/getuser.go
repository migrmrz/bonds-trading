package store

import "fmt"

func (dbs *DatabaseStorage) GetUser(username string) (User, error) {
	var users []User

	err := dbs.statementGetUser.Select(&users, username)
	if err != nil {
		return User{}, fmt.Errorf("error getting user data: %s", err.Error())
	}

	if len(users) == 0 {
		return User{}, fmt.Errorf("error: user not found")
	}

	return users[0], nil

}
