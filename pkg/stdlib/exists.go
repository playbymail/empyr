// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package stdlib

import (
	"errors"
	"io/fs"
	"os"
)

// IsExists returns true if the path exists.
func IsExists(path string) bool {
	return isExists(os.Stat(path))
}

// IsExistsFS returns true if the path exists.
// Accepts filesystem as a parameter, so this will work with any filesystem.
func IsExistsFS(path string, filesystem fs.FS) bool {
	return isExists(fs.Stat(filesystem, path))
}

// IsDirExists returns true only if the path exists and is a directory.
func IsDirExists(path string) bool {
	sb, err := os.Stat(path)
	return err == nil && sb.IsDir()
}

// IsDirExistsFS returns true if the path exists and is a directory.
//
// Accepts filesystem as a parameter, so this will work with any filesystem.
func IsDirExistsFS(path string, filesystem fs.FS) bool {
	sb, err := fs.Stat(filesystem, path)
	return err == nil && sb.IsDir()
}

// IsFileExists returns true only if the path exists and is a regular file.
// Returns false with no error if the path does not exist or is not a regular file.
// Returns false and an error if stat on the path fails for any other reason.
func IsFileExists(path string) bool {
	sb, err := os.Stat(path)
	return err == nil && sb.Mode().IsRegular()
}

// IsFileExistsFS returns true if the path exists and is a regular file.
//
// Accepts filesystem as a parameter, so this will work with any filesystem.
func IsFileExistsFS(path string, filesystem fs.FS) bool {
	sb, err := fs.Stat(filesystem, path)
	return err == nil && sb.Mode().IsRegular()
}

// isExists accepts the results of calling stat on a path.
// Returns true if the path exists.
// It is a helper function for IsExists and IsExistsFS.
func isExists(sb fs.FileInfo, err error) bool {
	return err == nil
}

// Stat accepts a path and returns the results of calling stat on the path.
// It handles the trivial case of the path not existing by returning nil, nil.
func Stat(path string) (fs.FileInfo, error) {
	sb, err := os.Stat(path)
	if err != nil {
		sb = nil
		if errors.Is(err, fs.ErrNotExist) {
			err = nil
		}
	}
	return sb, nil
}

// FSStat accepts a filesystem and path and returns the results of calling stat on the path.
// It handles the trivial case of the path not existing by returning nil, nil.
func FSStat(filesystem fs.FS, path string) (fs.FileInfo, error) {
	sb, err := fs.Stat(filesystem, path)
	if err != nil {
		sb = nil
		if errors.Is(err, fs.ErrNotExist) {
			err = nil
		}
	}
	return sb, nil
}

// Remove accepts a path and removes it if it is a regular file.
// Returns nil if the path does not exist.
// Returns an error if the path exists or there is an error removing it.
func Remove(path string) error {
	if !IsExists(path) {
		return nil
	} else if IsDirExists(path) {
		return errors.New("path is a directory")
	} else if !IsFileExists(path) {
		return errors.New("path is not a regular file")
	}
	return os.Remove(path)
}
