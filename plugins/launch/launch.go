package launch

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/eissar/nest/eagle/api"
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

	// resolve path
	if err := func() error {
		pth, _ := filepath.Abs(v)
		if _, err := os.Stat(pth); err != nil {
			if !errors.Is(err, os.ErrNotExist) { // error is not ErrNotExist
				return fmt.Errorf("unknown error testing path %s: err=%s", pth, err.Error())
			}
		} else {
			is_path = true
			v = pth // reassign to absolute filepath
		}
		return nil
	}(); err != nil {
		return err
	}

	if !is_path {
		fmt.Printf("api.IsValidItemID(v): %v\n", api.IsValidItemID(v))
		if api.IsValidItemID(v) {
			is_id = true
		}
	}

	if is_path {
		return revealFilePlatformSpecific(v)
	}
	if is_id {
		return openEagleItem(v)
	}

	return nil
}
