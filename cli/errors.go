// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package cli

// Error implements Cheney's constant error idiom
type Error string

// Error implements the error interface.
func (e Error) Error() string {
	return string(e)
}

const (
	ErrDeleteFailed      = Error("delete failed")
	ErrFileExists        = Error("file exists")
	ErrEnvPrefixBadChars = Error("environment prefix contains invalid characters")
	ErrEnvPrefixEmpty    = Error("environment prefix is empty")
	ErrEnvPrefixPrefix   = Error("environment prefix must start with an upper-case letter")
	ErrEnvPrefixSuffix   = Error("environment prefix must end with an underscore")
	ErrInvalidEnvPrefix  = Error("invalid environment prefix")
	ErrInvalidVersion    = Error("invalid version")
	ErrNotImplemented    = Error("not implemented")
)
