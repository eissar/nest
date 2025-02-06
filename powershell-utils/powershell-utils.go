package powershellutils

import (
	"encoding/json"
	"fmt"
	"os/exec"
)

func RunPwshCmd(ScriptPath string) []interface{} {
	cmd := exec.Command("pwsh.exe", "-NoProfile", "-c", ScriptPath)
	// Capture the output
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error executing command:", err)
		// You might want to handle specific exit codes differently:
		if exitError, ok := err.(*exec.ExitError); ok {
			fmt.Printf("Exit code: %d\n", exitError.ExitCode())
			fmt.Printf("Stderr: %s\n", exitError.Stderr)

			// Example: Check if PowerShell command didn't find any processes
			if exitError.ExitCode() == 1 { // Exit code 1 is returned when PowerShell doesn't error but the script returns $null or nothing (e.g. no processes found)
				fmt.Println("No Waterfox processes found.")
				panic(err)
			}
		}
	}
	var messageData []interface{}
	if err := json.Unmarshal(output, &messageData); err != nil {
		fmt.Printf("[ERROR]", err)
	}

	return messageData
}
func ExecPwshCmd(Cmd string) {
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
				fmt.Println("Null or empty output.")
				panic(err)
			}
		}
	}
}
