// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package empyr

import "github.com/maloquacious/cerrors"

// all visible errors raised by this package
const (
	ErrDistanceExceedsCapacity = cerrors.Error("distance exceeds capacity")
	ErrInsufficientFuel        = cerrors.Error("insufficient fuel")
	ErrMassExceedsCapacity     = cerrors.Error("mass exceeds capacity")
)
