// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package ratelimiter

type Limiter struct { }

func (l *Limiter) Allow(key string, limit int) bool {
	return true
}
