// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package store

// Error implements Cheney's constant error idiom
type Error string

// Error implements the error interface.
func (e Error) Error() string {
	return string(e)
}

const (
	ErrAlreadyExists       = Error("already exists")
	ErrCreateSchema        = Error("create schema")
	ErrForeignKeysDisabled = Error("foreign keys are disabled")
	ErrInvalidPath         = Error("invalid path")
	ErrNotExist            = Error("not exist")
	ErrNotFound            = Error("not found")
	ErrNotImplemented      = Error("not implemented")
	ErrNotOpen             = Error("not open")
	ErrNotUnique           = Error("not unique")
	ErrNotValid            = Error("not valid")
	ErrNotWritable         = Error("not writable")
	ErrOpen                = Error("open")
	ErrPragmaReturnedNil   = Error("pragma returned nil")
	ErrReadOnly            = Error("read only")
	ErrUnknown             = Error("unknown")
	ErrUnsupported         = Error("unsupported")
	ErrWriteOnly           = Error("write only")
)
