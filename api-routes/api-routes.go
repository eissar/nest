package apiroutes

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"
	pwsh "web-dashboard/powershell-utils"
	. "web-dashboard/types"

	"github.com/labstack/echo/v4"
)

var editor = "C:/Program Files/Neovim/bin/nvim.exe"

func Ping(c echo.Context) error {
	return c.String(http.StatusOK, "OK")
}
func Edit(c echo.Context) error {
	start := time.Now()
	leaf := c.QueryParam("fileName") // Get-ChildItem "$env:CLOUD_DIR/Catalog/*.md"
	fp := path.Join(os.Getenv("CLOUD_DIR"), "Catalog", leaf)
	fmt.Println(editor, fp)
	//fp = "C:/Users/eshaa/draft.lua"
	open_cmd := fmt.Sprintf("wt.exe -d $env:CLOUD_DIR nvim %v", fp)
	pwsh.ExecPwshCmd(open_cmd)
	fmt.Println("[Debug] (", c.Path(), ") request elapsed time:", time.Since(start))
	return c.String(http.StatusOK, "success")
}

func RecentEagleItems(c echo.Context) error {
	//  TODO: IMPROVE SPEED (791.9206ms)
	start := time.Now()
	a := pwsh.RunPwshCmd("./recentEagleItems.ps1")
	fmt.Println("[Debug] (", c.Path(), ") request elapsed time:", time.Since(start))
	return c.JSON(200, a)
}

func RecentNotes(c echo.Context) error {
	start := time.Now()
	a := pwsh.RunPwshCmd("./recentNotes.ps1")
	fmt.Println("[Debug] (", c.Path(), ") request elapsed time:", time.Since(start))
	return c.JSON(200, a)
}

func NumTabs(c echo.Context) error {
	start := time.Now()
	a := pwsh.RunPwshCmd("./waterfoxTabs.ps1")
	fmt.Println("[Debug] (", c.Path(), ") request elapsed time:", time.Since(start))
	return c.JSON(200, a)
}

// returns: array of window structs.
func GetEnumerateWindows() []Window {
	cmd := exec.Command("C:/Users/eshaa/Dropbox/Code/cs/enumerateWindowsExe/aot/enumerateWindows.exe")
	var stdoutBuf, stderrBuf bytes.Buffer
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf

	err := cmd.Run()
	if err != nil {
		// Capture stderr if the command fails
		if _, ok := err.(*exec.ExitError); ok {
			// fmt.Println("executable exited with error: %v, status code: %d, stderr: %s", exitError, exitError.ExitCode(), stderrBuf.String())
			return nil
		}
		panic(err)
	}
	// finish running exec

	// just do everything without thinking too hard since it's not that much data.

	unparsed := stdoutBuf.String()
	parsed := strings.ReplaceAll(unparsed, `"`, "'")

	fmt.Println()
	reader := csv.NewReader(strings.NewReader(parsed))

	reader.LazyQuotes = true    // Handle double quotes within fields
	reader.Comma = '|'          // delimiter pipe '|'
	reader.FieldsPerRecord = -1 // -1 for any len

	var data []Window

	_, err = reader.Read() // Skip the header row if exists
	if err != nil {
		panic(err)
	}

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break // End of file
		}
		if err != nil {
			panic(err)
		}

		var item Window
		if len(record) == 3 {
			item = Window{
				Handle:    record[0],
				Title:     record[1],
				ProcessId: record[2],
			}
			data = append(data, item)
		}
	}
	return data
}
func EnumWindows(c echo.Context) error {
	a := GetEnumerateWindows()
	return c.JSON(200, a)
}
