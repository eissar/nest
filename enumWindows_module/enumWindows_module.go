package enumwindowsmodule

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	. "web-dashboard/types"
)

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

func PopulateEnumerateWindows(c echo.Context, templateName string) interface{} {
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
	firstParameter := c.QueryParam("first")
	//firstParameter := c.Param("first")
	dataLen := len(data)
	if dataLen == 0 {
		return nil
	}
	if firstParameter != "" {
		f, err := strconv.ParseUint(firstParameter, 10, 64)
		if err != nil {
			msg := fmt.Sprintf(
				"Error populating template %s . parameter `first` was not a valid integer. Error: %s",
				templateName,
				err.Error(),
			)
			fmt.Println("[ERROR]", msg)
			e := fmt.Sprintf("<p>%s</p>", msg)
			c.String(http.StatusBadRequest, e)
			return nil
		}
		// return the first `f` or the whole array,
		// whichever is bigger
		if dataLen < int(f) {
			f = uint64(dataLen)
		}
		// TODO: create architecture for some kind of response annotation somehow.
		return data[0:f]
	}
	return data
}

func Enum(c echo.Context) error {
	a := GetEnumerateWindows()
	return c.JSON(200, a)
}

func RegisterRoutesFromGroup(g *echo.Group) {}
func RegisterRootRoutes(server *echo.Echo) {
	server.GET("/api/windows", Enum)
}
