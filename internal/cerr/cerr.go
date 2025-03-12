// Copyright (c) 2025 Michael D Henderson. All rights reserved.

// Package cerr implements Cheney's constant error idiom.
package cerr

// Error defines a constant error.
type Error string

// Error implements the error interface.
func (e Error) Error() string {
	return string(e)
}
