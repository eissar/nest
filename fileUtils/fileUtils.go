package fileUtils

import (
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
)

var ( // copied from os
	// ErrInvalid indicates an invalid argument.
	// Methods on File will return this error when the receiver is nil.
	ErrInvalid = fs.ErrInvalid // "invalid argument"

	ErrPermission = fs.ErrPermission // "permission denied"
	ErrExist      = fs.ErrExist      // "file already exists"
	ErrNotExist   = fs.ErrNotExist   // "file does not exist"
	ErrClosed     = fs.ErrClosed     // "file already closed"
)

// returns absolute filepath
// or log.Fatalf()
func mustFilepathAbs(s string) string {
	fp, err := filepath.Abs(s)
	if err != nil {
		log.Fatalf("pathexists: %w,", err)
	}
	return fp
}

func TestPath(pth string) error {
	_, err := os.Stat(pth)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("testpath: import error, file `%s` does not exist! err: %w", pth, err)
		}
		return fmt.Errorf("testpath: import error, unknown error checking if path `%s` exists. err: %w", pth, err)
	}
	return nil
}
