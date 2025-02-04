package main

import (
	_ "bytes"
	_ "encoding/csv"
	"fmt"
	_ "io"
	"os"
	_ "os/exec"
	"path/filepath"
	_ "strings"
	_ "web-dashboard/types"

	"github.com/labstack/echo/v4"
)

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
