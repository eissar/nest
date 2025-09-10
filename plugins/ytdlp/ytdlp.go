package ytdlp

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/eissar/nest/api"
	"github.com/eissar/nest/config"
	// "github.com/spf13/cobra"
)

// is ytdlp
func AssertAvailable() {
	cmd := exec.Command("yt-dlp.exe", "--version")
	// Capture the output
	output, err := cmd.CombinedOutput()
	if err != nil {
		panic(err)
	}

	fmt.Printf("output: %s\n", string(output))
}

func Get(target string) {
	tempDir, err := os.MkdirTemp("", "nest-")
	if err != nil {
		panic(fmt.Sprintf("COULD NOT CREATE TEMP FILE\n %s", err))
	}
	// defer fmt.Printf("Downloaded to: %s\n", tempDir)

	outputPath := filepath.Join(tempDir, "%(title)s.%(ext)s")
	cmd := exec.Command("yt-dlp.exe", "-o", outputPath, target)
	_, err = cmd.CombinedOutput()
	if err != nil {
		panic(err)
	}

	// fmt.Printf("\n===output===\n%v\n", output)

	fmt.Printf("Downloaded to %s\n", tempDir)

	entries, err := os.ReadDir(tempDir)
	if err != nil {
		panic(err)
	}

	var containsSt = false

	// just make sure something is in here
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		containsSt = true
		break
	}
	// so with add from path we can actually just
	// direct it at a path and it will add everything automatically
	// so we can use this hidden behavior to make this pretty simple
	if containsSt {
		opts := api.ItemAddFromPathOptions{Path: tempDir}
		opts.Validate()

		cfg := config.GetConfig()

		api.ItemAddFromPath(cfg.BaseURL(), opts)
	} else {
		panic(fmt.Sprintf("tempdir at %s contains no items!", tempDir))
	}
}
