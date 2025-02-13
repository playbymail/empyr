// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package main

import "errors"

// ErrBadRequest is used when creating an object and something is wrong with the request.
var ErrBadRequest = errors.New("bad request")

// ErrDuplicate is used when the object is not unique.
var ErrDuplicate = errors.New("duplicate")

// ErrDuplicateAddress is used when the e-mail address is not unique.
var ErrDuplicateAddress = errors.New("duplicate e-mail address")

var ErrDuplicatePlayer = errors.New("duplicate player")

// ErrDuplicateUserName is used when the user name is not unique.
var ErrDuplicateUserName = errors.New("duplicate user name")

// ErrNoData is used when a game could not be found.
var ErrNoData = errors.New("no data found")

// ErrUnknown is used when the error source is unknown.
var ErrUnknown = errors.New("internal server error")
