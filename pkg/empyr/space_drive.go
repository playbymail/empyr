// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package empyr

type OrbitalDrive = SpaceDrive

type SpaceDrive struct {
	TechLevel int
	Quantity  int
	Mass      float64 // total mass, in tonnes
	Volume    float64 // total volume, in cubic meters
}
