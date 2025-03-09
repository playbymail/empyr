// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package engine

// Error implements Cheney's constant error idiom
type Error string

// Error implements the error interface.
func (e Error) Error() string {
	return string(e)
}
