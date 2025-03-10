// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package engine

import (
	"errors"
	"fmt"
	"regexp"
)

const (
	ErrInvalidCode        = Error("invalid code")
	ErrInvalidDescription = Error("invalid description")
	ErrInvalidHandle      = Error("invalid handle")
	ErrInvalidName        = Error("invalid name")
	ErrMissingHandle      = Error("missing handle")
)

// IsValidCode validates whether a game code meets specific formatting requirements.
// It checks that the code is between 3 and 5 characters long and starts with uppercase letters,
// optionally followed by numbers. Returns true if the code is valid, otherwise returns false
// with an error describing the validation failure.
func IsValidCode(code string) (bool, error) {
	var validGameCode = regexp.MustCompile(`^[A-Z]+[0-9]*$`)
	if !(3 <= len(code) && len(code) <= 5) {
		return false, errors.Join(fmt.Errorf("code must be between 3...5 characters"), ErrInvalidCode)
	} else if !validGameCode.MatchString(code) {
		return false, errors.Join(fmt.Errorf("code must start with uppercase letters, optionally followed by numbers"), ErrInvalidCode)
	}
	return true, nil
}

// IsValidDescription validates whether a description meets the same character requirements as a name.
// It delegates the validation to IsValidName, checking that the description contains only allowed
// characters (alphanumeric, underscore, hyphen, dot, and space) and does not contain any special
// escape sequences or unusual characters.
//
// Returns true if the description is valid, otherwise returns false with an error describing
// the validation failure. The error indicates that the description must not contain special characters.
func IsValidDescription(descr string) (bool, error) {
	return IsValidName(descr)
}

// IsValidHandle validates whether a handle meets specific formatting requirements.
func IsValidHandle(handle string) (bool, error) {
	var validHandle = regexp.MustCompile(`^[a-z][a-zA-Z0-9_]+$`)
	if !validHandle.MatchString(handle) {
		return false, errors.Join(fmt.Errorf("handle must contain only alphanumeric characters, numbers, and underscores"), ErrInvalidHandle)
	}
	return true, nil
}

// IsValidName validates whether a name meets specific character requirements.
// It checks that the name contains only allowed characters (alphanumeric, underscore, hyphen, dot, and space)
// and does not contain any special escape sequences or unusual characters.
//
// The function performs two checks:
// 1. Each character in the name must be in the predefined goodNameBytes set
// 2. The name must not be considered "weird" (containing special characters or escape sequences)
//
// Returns true if the name is valid, otherwise returns false with an error describing the validation failure.
// The error indicates that the name must not contain special characters.
func IsValidName(name string) (bool, error) {
	for _, ch := range []byte(name) {
		if !goodNameBytes[ch] {
			//log.Printf("name: %q: %q\n", name, string(ch))
			return false, errors.Join(fmt.Errorf("name must not contain special characters"), ErrInvalidName)
		}
	}
	if isWeird(name) {
		//log.Printf("name: %q: isWeird\n", name)
		return false, errors.Join(fmt.Errorf("name must not contain special characters"), ErrInvalidName)
	}
	return true, nil
}

var (
	goodNameBytes [256]bool
)

func init() {
	for _, ch := range []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_-. ") {
		goodNameBytes[ch] = true
	}
}

// isWeird checks if a string contains any special characters or escape sequences
// by comparing the raw string with its quoted representation. It's a clever way
// to detect special characters by leveraging Go's string quoting behavior.
//
// Returns true if the string contains special characters, false otherwise.
// Example: isWeird("hello") returns false, isWeird("hello\n") returns true
func isWeird(s string) bool {
	return `"`+s+`"` != fmt.Sprintf("%q", s)
}
