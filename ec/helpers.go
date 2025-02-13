// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package ec

import (
	"encoding/json"
	"os"
	"path/filepath"
)

func fromjson(path, name string, data any) error {
	buffer, err := os.ReadFile(filepath.Clean(filepath.Join(path, name+".json")))
	if err != nil {
		return err
	}
	return json.Unmarshal(buffer, data)
}

func tojson(path, name string, data any) error {
	buffer, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Clean(filepath.Join(path, name+".json")), buffer, 0644)
}
