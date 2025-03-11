// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package flash

import "net/http"

// GetStore returns the flash store for the current request.
func GetStore(r *http.Request) *Store {
	return nil
}

// Store is a flash store keyed to the current request.
type Store struct{}

func (s *Store) Set(key, value string) {}
func (s *Store) Get(key string) string { return "" }
