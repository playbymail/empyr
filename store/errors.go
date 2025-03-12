// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package store

import "github.com/playbymail/empyr/internal/cerr"

const (
	ErrAlreadyExists       = cerr.Error("already exists")
	ErrCreateSchema        = cerr.Error("create schema")
	ErrForeignKeysDisabled = cerr.Error("foreign keys are disabled")
	ErrInvalidPath         = cerr.Error("invalid path")
	ErrNotExist            = cerr.Error("not exist")
	ErrNotFound            = cerr.Error("not found")
	ErrNotImplemented      = cerr.Error("not implemented")
	ErrNotOpen             = cerr.Error("not open")
	ErrNotUnique           = cerr.Error("not unique")
	ErrNotValid            = cerr.Error("not valid")
	ErrNotWritable         = cerr.Error("not writable")
	ErrOpen                = cerr.Error("open")
	ErrPragmaReturnedNil   = cerr.Error("pragma returned nil")
	ErrReadOnly            = cerr.Error("read only")
	ErrUnknown             = cerr.Error("unknown")
	ErrUnsupported         = cerr.Error("unsupported")
	ErrWriteOnly           = cerr.Error("write only")
)
