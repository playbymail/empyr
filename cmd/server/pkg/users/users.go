// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package users

import (
	"errors"
	"github.com/google/uuid"
	"time"
)

// User defines the properties of a user.
type User struct {
	ID      string    `json:"id"`
	Email   string    `json:"email"`
	Name    string    `json:"name"`
	Created time.Time `json:"created"`
}

// Users is a map of user data
type Users struct {
	data map[string]*User
}

// ErrDuplicateAddress is used when the e-mail address is not unique.
var ErrDuplicateAddress = errors.New("duplicate e-mail address")

// ErrDuplicateName is used when the user name is not unique.
var ErrDuplicateName = errors.New("duplicate user name")

func New() *Users {
	return &Users{data: make(map[string]*User)}
}

func (u *Users) All() []*User {
	var list []*User = []*User{}
	for _, user := range u.data {
		list = append(list, user)
	}
	return list
}

func (u *Users) Create(id, name, email string) (*User, error) {
	for _, user := range u.data {
		if user.Email == email {
			return user, ErrDuplicateAddress
		} else if user.Name == name {
			return user, ErrDuplicateName
		}
	}

	if id == "" { // todo: this is just for development
		id = uuid.New().String()
	}

	user := &User{
		ID:      id,
		Email:   email,
		Name:    name,
		Created: time.Now(),
	}
	u.data[user.ID] = user

	return user, nil
}

func (u *Users) Filter(fn func(*User) bool) []*User {
	var list []*User = []*User{}
	for _, user := range u.data {
		if fn(user) {
			list = append(list, user)
		}
	}
	return list
}

func (u *Users) ID(id string) *User {
	return u.data[id]
}
