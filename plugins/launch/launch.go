package launch

import (
	"fmt"
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
