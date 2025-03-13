// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package cli

import "github.com/playbymail/empyr/internal/cerr"

const (
	ErrDeleteFailed      = cerr.Error("delete failed")
	ErrFileExists        = cerr.Error("file exists")
	ErrEnvFlagInvalid    = cerr.Error("environment flag is invalid")
	ErrEnvFlagNotSet     = cerr.Error("environment flag not set")
	ErrEnvPrefixBadChars = cerr.Error("environment prefix contains invalid characters")
	ErrEnvPrefixEmpty    = cerr.Error("environment prefix is empty")
	ErrEnvPrefixPrefix   = cerr.Error("environment prefix must start with an upper-case letter")
	ErrEnvPrefixSuffix   = cerr.Error("environment prefix must end with an underscore")
	ErrInvalidEnvPrefix  = cerr.Error("invalid environment prefix")
	ErrInvalidVersion    = cerr.Error("invalid version")
	ErrNotImplemented    = cerr.Error("not implemented")
)
