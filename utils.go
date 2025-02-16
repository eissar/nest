package main

import (
	_ "bytes"
	_ "encoding/csv"
	"fmt"
	_ "github.com/eissar/nest/types"
	_ "io"
	"os"
	"os/exec"
	_ "os/exec"
	"path/filepath"
	"runtime"
	_ "strings"

	"github.com/labstack/echo/v4"
)

// utils_module
// defines `hanging` handler functions
// which are not tightly coupled to other functionality,
// state, or whose functionality are entirely self-contained

func readFile(p string) []byte {
	fp, err := filepath.Abs(p)
	if err != nil {
		panic(err)
	}
	filebytes, err := os.ReadFile(fp)
	if err != nil {
		panic(err)
	}
	return filebytes
}

func PrintSiteMap(server *echo.Echo) {
	fmt.Println("server available routes:")
	for _, x := range server.Routes() {
		fmt.Println(x.Name, x.Path)
	}
}

func openURI(uri string) error {
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
