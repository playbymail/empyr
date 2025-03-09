// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package cli

import (
	"errors"
	"github.com/mdhender/semver"
	"strings"
)

type Option func() error

// WithEnvPrefix sets the prefix for environment variables.
// Ensures that the prefix:
//  1. is not empty
//  2. is upper-case
//  3. starts with an upper-case letter
//  4. contains no special characters
func WithEnvPrefix(pfx string) Option {
	const (
		validChars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"
	)
	return func() error {
		if pfx == "" {
			return errors.Join(ErrEnvPrefixEmpty, ErrInvalidEnvPrefix)
		} else if !('A' <= pfx[0] && pfx[0] <= 'Z') {
			return errors.Join(ErrEnvPrefixPrefix, ErrInvalidEnvPrefix)
		} else if containsInvalidChars(pfx, validChars) {
			return errors.Join(ErrEnvPrefixBadChars, ErrInvalidEnvPrefix)
		}
		env.Env.Prefix = pfx
		return nil
	}
}

func WithVersion(version semver.Version) Option {
	return func() error {
		env.Version = version
		return nil
	}
}

// containsInvalidChars returns true if s contains any character not in validChars
func containsInvalidChars(s, validChars string) bool {
	for _, c := range s {
		if !strings.ContainsRune(validChars, c) {
			return true
		}
	}
	return false
}
