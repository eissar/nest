package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/eissar/nest/config"
	"github.com/eissar/nest/core"
	cmd "github.com/eissar/nest/core/command-line"
)

// TODO:
// - [ ] ? replace flag with github.com/spf13/cobra
// split into start.exe / nest.exe

func main() {
	addCmd := flag.NewFlagSet("add", flag.ExitOnError)
	addPath := addCmd.String("file", "", "filepath that will be added to eagle")

	listCmd := flag.NewFlagSet("list", flag.ExitOnError)
	listLimit := listCmd.Int("limit", 5, "number of items to retrieve")

	revealCmd := flag.NewFlagSet("reveal", flag.ExitOnError)
	revealPath := revealCmd.String("target", "", "filepath or item id to reveal")

	switchCmd := flag.NewFlagSet("switch", flag.ExitOnError)
	switchName := switchCmd.String("name", "", "name of library to switch to.")
	//revealCmd := flag.NewFlagSet("reveal", flag.ExitOnError)

	help := flag.Bool("help", false, "print help information")
	start := flag.Bool("start", false, "run the utility server")
	//debug := flag.Bool("debug", true, "shows additional information in the console while running. (does nothing)")
	stop := flag.Bool("stop", false, "stop the utility server")
	flag.Parse()

	if *help || len(os.Args) < 2 {
		fmt.Println("expected flag or subcommand.")
		flag.PrintDefaults()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "add":
		cfg := config.GetConfig()

		addCmd.Parse(os.Args[2:])
		cmd.Add(cfg, addPath)
		os.Exit(0)
	case "list":
		cfg := config.GetConfig()

		listCmd.Parse(os.Args[2:])
		cmd.List(cfg, listLimit)
		os.Exit(0)
	case "reveal":
		cfg := config.GetConfig()

		revealCmd.Parse(os.Args[2:])
		cmd.Reveal(cfg, revealPath)
		os.Exit(0)
	case "switch":
		cfg := config.GetConfig()
		switchCmd.Parse(os.Args[2:])
		if *switchName != "" {
			cmd.Switch(cfg, *switchName)
		}
		if len(os.Args) < 3 {
			log.Fatalf("must pass flag -name")
			flag.PrintDefaults()
			os.Exit(1)
		}
		cmd.Switch(cfg, os.Args[2])
		os.Exit(0)
	}

	//if *debug { } /* pwsh.ExecPwshCmd("./powershell-utils/openUrl.ps1 -Uri 'http://localhost:1323/app/notes'") */

	if *start {
		core.Start() //blocking
	}

	if *stop {
		err := cmd.Shutdown(config.GetConfig())
		if err != nil {
			fmt.Printf("stop: %s", err.Error())
		}
		os.Exit(0)
	}
}
