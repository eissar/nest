package cmd

import (
	"fmt"
	"os"
	"path/filepath"
)

type FileDir string

func (d *FileDir) String() string { return string(*d) }

func (d *FileDir) Set(value string) error {
	// Resolve symlinks (optional, but makes error clearer)
	abs, err := filepath.Abs(value)
	if err != nil {
		return fmt.Errorf("cannot resolve absolute path %q: %w", value, err)
	}

	info, err := os.Stat(abs)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("directory %q does not exist", abs)
		}
		return fmt.Errorf("cannot stat %q: %w", abs, err)
	}
	if !info.IsDir() {
		return fmt.Errorf("path %q is not a directory", abs)
	}

	// check for write permission: if info.Mode().Perm()&(0200) == 0 { … }

	*d = FileDir(abs) // store the *absolute* (or original) path
	return nil
}

func (d *FileDir) Type() string { return "directory" }

// we validate the filepath.
// this is just for flags so is fine
// but don't instantiate a bunch of these
type FilePath string

func (f *FilePath) String() string { return string(*f) }

func (f *FilePath) Set(value string) error {
	// Resolve symlinks (but makes error clearer)
	abs, err := filepath.Abs(value)
	if err != nil {
		return fmt.Errorf("cannot resolve absolute path %q: %w", value, err)
	}

	info, err := os.Stat(abs)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("file %q does not exist", abs)
		}
		return fmt.Errorf("cannot stat %q: %w", abs, err)
	}
	if info.IsDir() {
		return fmt.Errorf("path %q is a directory, expected a file", abs)
	}

	*f = FilePath(abs) // store the absolute (or original) path
	return nil
}

// Type is used by pflag to describe the flag in help output.
func (f *FilePath) Type() string { return "file" }
