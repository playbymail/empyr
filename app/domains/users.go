// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package domains

import "errors"

type User struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

type UserRepository interface {
	Save(user User) error
}

type UserService struct {
	Repo UserRepository
}

func (s *UserService) CreateUser(username, email string) (User, error) {
	if username == "" || email == "" {
		return User{}, errors.New("invalid input")
	}

	user := User{Username: username, Email: email}
	err := s.Repo.Save(user)
	return user, err
}
