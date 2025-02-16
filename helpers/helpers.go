package helpers

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
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

// panics if path does not exist.
func AssertPathExists(pth string) {
	_, err := os.Stat(pth) // More efficient than opening and closing
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			log.Fatalf("import: import error, file `%s` does not exist! err: %s", pth, err)
		}
		log.Fatalf("import: import error, unknown error checking if path `%s` exists. err: %s", pth, err)
	}
}
