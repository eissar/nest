package pwsh

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

//type ExecErr error

// TODO: make this not stupid
// parameter script: filename of script in pwsh module.
// (e.g., scriptName.ps1)
// if there is an error, returns interface
// []error with an error in interface...
func RunScript(script string) []interface{} {
	pth := filepath.Join("./modules/pwsh", script)
	pth, err := filepath.Abs(pth)
	if err != nil {
		return []interface{}{err}
	}
	_, err = os.Stat(pth)
	if err != nil {
		return []interface{}{err}
	}
	cmd := exec.Command("pwsh.exe", "-NoProfile", "-c", pth)
	// Capture the output
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error executing command:", err)
		if exitError, ok := err.(*exec.ExitError); ok {
			return []interface{}{
				fmt.Errorf(
					"RunScript: error executing command %s exitcode=%s stderr=%s",
					exitError.Error(),
					exitError.ExitCode(),
					exitError.Stderr,
				),
			}
		}
		return []interface{}{
			fmt.Errorf("RunScript: err=%w", err),
		}
	}
	var messageData []interface{}
	if err := json.Unmarshal(output, &messageData); err != nil {
		if err != nil {
			return []interface{}{err}
		}
	}

	return messageData
}
func Exec(Cmd string) {
	cmd := exec.Command("pwsh.exe", "-NoProfile", "-c", Cmd)
	// Capture the output
	_, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error executing command:", err)
		// You might want to handle specific exit codes differently:
		if exitError, ok := err.(*exec.ExitError); ok {
			fmt.Printf("Exit code: %d\n", exitError.ExitCode())
			fmt.Printf("Stderr: %s\n", exitError.Stderr)

			// Example: Check if PowerShell command didn't find any processes
			if exitError.ExitCode() == 1 { // Exit code 1 is returned when PowerShell doesn't error but the script returns $null or nothing (e.g. no processes found)
				fmt.Println("[ERROR] <> Null or empty output. Error", exitError.Error())
			}
		}
	}
}

//func RegisterGroupRoutes(server *echo.Echo) {
//	g.GET("/api/recentEagleItems", RecentEagleItems)
//}
