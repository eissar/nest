package launch

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/eissar/eagle-go"
)

func OpenURI(uri string) error {
	var cmd *exec.Cmd

	fmt.Println("[LOG] <openUri> opening...", uri)
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", uri)
	case "darwin":
		cmd = exec.Command("open", uri)
	default: // Linux and other Unix-like systems
		cmd = exec.Command("xdg-open", uri)
	}

	return cmd.Run()
}

func Open(file string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("explorer", file)
	case "darwin":
		cmd = exec.Command("open", file)
	default: // Linux and other Unix-like systems
		cmd = exec.Command("xdg-open", file)
	}

	cmd.Run()
	return nil
}

// param v string: filepath or item id
// TODO: reveal fails if filepath has spaces.
// TODO: fix logic to make itemID case retrieve filepath and reveal in explorer
func Reveal(v string) error {
	revealFilePlatformSpecific := func(filePath string) error {
		switch runtime.GOOS {
		case "windows":
			//cmd := exec.Command(`explorer`, `/select,`, filePath)
			//cmd.Run()
			exec.Command(`explorer`, `/select,`, filePath).Run()
		case "darwin":
			return fmt.Errorf("not implemented: %w", errors.ErrUnsupported)
		default: // Linux and other Unix-like systems
			return fmt.Errorf("not implemented: %w", errors.ErrUnsupported)
		}
		return nil
	}

	openEagleItem := func(id string) error {
		OpenURI("eagle://item/" + id)
		return nil
	}

	is_path := false
	is_id := false

	pth, err := filepath.Abs(v)
	if err != nil {
		// The original code ignored this potential error
		return fmt.Errorf("failed to resolve absolute path for %q: %w", v, err)
	}
	if a, err := os.Lstat(pth); err == nil {
		is_path = true
	} else if !errors.Is(err, os.ErrNotExist) {
		// An unexpected error occurred (e.g., permission denied)
		return fmt.Errorf("error checking path %q: %w", pth, err)
	} else if errors.Is(err, os.ErrNotExist) {
		// If err is os.ErrNotExist, do nothing and continue
		// fmt.Println("exists?", e, err)

		fmt.Println("filepath:", pth, "error", err.Error())
		fmt.Println(a)
		fmt.Println("file does not exist?")
	}

	if !is_path {
		fmt.Printf("api.IsValidItemID(%v): %v\n", pth, eagle.IsValidItemID(v))
		if eagle.IsValidItemID(v) {
			is_id = true
		}
	}

	if is_path {
		return revealFilePlatformSpecific(pth)
	}
	if is_id {
		return openEagleItem(v)
	}

	return nil
}
