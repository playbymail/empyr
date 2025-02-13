// Copyright (c) 2025 Michael D Henderson. All rights reserved.

// Package clean implements functions to check for bad characters in
// codes, names, and descriptions.
package clean

import (
	"errors"
	"fmt"
	"regexp"
)

func IsValidCode(code string) (bool, error) {
	var validGameCode = regexp.MustCompile(`^[A-Z]+[0-9]*$`)
	if !(3 <= len(code) && len(code) <= 5) {
		return false, errors.Join(fmt.Errorf("code must be between 3...5 characters"), ErrInvalidCode)
	} else if !validGameCode.MatchString(code) {
		return false, errors.Join(fmt.Errorf("code must start with uppercase letters, optionally followed by numbers"), ErrInvalidCode)
	}
	return true, nil
}

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
