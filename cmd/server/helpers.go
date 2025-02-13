// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package main

import "math/rand"

// roll returns the sum of N D-sided dice.
func roll(prng *rand.Rand, n, d int) (sum int) {
	for i := 0; i < n; i++ {
		sum = sum + prng.Intn(d) + 1
	}
	return sum
}
