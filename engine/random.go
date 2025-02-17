// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package engine

import (
	"math"
	"math/rand/v2"
)

// this file implements functions to return random values.

// bellCurve returns a random value following a normal (Gaussian) distribution
// with the specified mean and standard deviation. The inputs are not validated,
// so callers must ensure they are appropriate for their use case.
//
// The standard deviation (stdDev) controls how spread out the values are:
// - 68% of values fall within ±1 stdDev of the mean
// - 95% of values fall within ±2 stdDev of the mean
// - 99.7% of values fall within ±3 stdDev of the mean
//
// For bounded ranges, using stdDev = (max-min)/6 ensures ~99.7% of values
// fall within the desired bounds. Values should be clamped to enforce limits.
//
// Example usage:
//
//	r := rand.New(rand.NewSource(uint64(time.Now().UnixNano())))
//
//	// Generate values clustered around 50 with most falling between 40-60
//	value := bellCurve(r, 50, 10)
func bellCurve(r *rand.Rand, mean, stdDev float64) float64 {
	return r.NormFloat64()*stdDev + mean
}

// normalRandInRange returns a random integer between min and max (inclusive) following
// a normal distribution centered between the values. The function handles reversed
// inputs and equal min/max values automatically.
//
// The returned values cluster around the midpoint between min and max, with fewer
// values appearing near the edges of the range, creating a natural distribution.
//
// Example usage:
//
//	r := rand.New(rand.NewSource(uint64(time.Now().UnixNano())))
//
//	// Generate deposit sizes that cluster around 50
//	typicalDeposit := normalRandInRange(r, 1, 99)
//
//	// Generate values for a rich region
//	richDeposit := normalRandInRange(r, 50, 99)
//
//	// Generate values for a poor region
//	poorDeposit := normalRandInRange(r, 1, 50)
func normalRandInRange(r *rand.Rand, min, max int) int {
	// ensure lower < upper
	var lower, upper float64
	if min < max {
		lower, upper = float64(min), float64(max)
	} else if min == max {
		lower, upper = float64(min), float64(max+1)
	} else {
		lower, upper = float64(max), float64(min)
	}

	// Calculate the mean and standard deviation.
	// Use 1/6 of the range as stdDev to keep ~99.7% of values within bounds
	mean := float64(lower+upper) / 2
	stdDev := float64(upper-lower) / 6

	// Rounding to nearest integer ensures the returned value is within bounds
	// and helps to preserve the natural distribution.
	value := math.Round(r.NormFloat64()*stdDev + mean)

	// Clamp to ensure we stay within bounds
	if value < lower {
		value = lower
	} else if value > upper {
		value = upper
	}

	return int(value)
}

// randomPoint returns scaled coordinates with a uniform volume distribution
// within a sphere of the given radius.
func randomPoint(radius float64) Point_t {
	// generate a random distance to ensure uniform distribution within the sphere
	d := math.Cbrt(rand.Float64()) // Cube root to ensure uniform distribution
	d = radius * d                 // Scale to the desired radius

	// generate random angles for spherical coordinates
	theta := rand.Float64() * 2 * math.Pi  // 0 to 2π
	phi := math.Acos(2*rand.Float64() - 1) // 0 to π

	// convert spherical coordinates to Cartesian coordinates
	return Point_t{
		X: int(math.Round(d * math.Sin(phi) * math.Cos(theta))),
		Y: int(math.Round(d * math.Sin(phi) * math.Sin(theta))),
		Z: int(math.Round(d * math.Cos(phi))),
	}
}

// randomCubePoint returns coordinates within a cube with side length of 31,
// centered at the origin and uniformly distributed.
func randomCubePoint() Point_t {
	return Point_t{
		X: rand.IntN(31) - 15, // -15 to 15
		Y: rand.IntN(31) - 15, // -15 to 15
		Z: rand.IntN(31) - 15, // -15 to 15
	}
}
