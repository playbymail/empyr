// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package clean

// Error implements Cheney's constant error idiom
type Error string

// Error implements the error interface.
func (e Error) Error() string {
	return string(e)
}

const (
	ErrInvalidCode = Error("invalid code")
	ErrInvalidName = Error("invalid name")
)
