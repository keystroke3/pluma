package pluma

import (
	"os"
)

type IsNotFileError struct {
	Path string
}

func (e *IsNotFileError) Error() string { return e.Path + " is not a file" }

// Checks if given path string p exists and is a file. Returns nil if both conditions are met.
// If it does not exist, os.DoesNotExist error is returned
// if it exists but is not a file, then IsNotFileError error is returned.
func isFile(p string) error {
	if p == "" {
		return &IsNotFileError{p}
	}
	stat, err := os.Stat(p)
	if err != nil && os.IsNotExist(err) {
		return err
	}
	if stat.IsDir() {
		return &IsNotFileError{p}
	}
	return nil
}
