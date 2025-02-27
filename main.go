package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/eissar/nest/config"
	"github.com/eissar/nest/core"
	"github.com/eissar/nest/eagle/api"
)

// globals
var debug = false

func main() {
	addCmd := flag.NewFlagSet("add", flag.ExitOnError)
	addPath := addCmd.String("file", "", "filepath that will be added to eagle.")

	help := flag.Bool("help", false, "print help information")
	serve := flag.Bool("serve", false, "run the utility server")
	//debug := flag.Bool("debug", true, "shows additional information in the console while running. (does nothing)")
	flag.Parse()

	if *help || len(os.Args) < 2 {
		fmt.Println("expected flag or subcommand.")
		flag.PrintDefaults()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "add":
		addCmd.Parse(os.Args[2:])
		fmt.Println("path:", *addPath)
		// resolve *addPath
		if len(*addPath) == 0 {
			log.Fatalf("[ERROR] add: flag `-file` is required.")
		}
		cfg := config.GetConfig()

		err := api.AddItemFromPath(cfg.FmtURL(), *addPath)
		if err != nil {
			log.Fatalf("[ERROR] while adding eagle item: err=%s", err.Error())
		}
		os.Exit(0)
	}

	//if *debug { } /* pwsh.ExecPwshCmd("./powershell-utils/openUrl.ps1 -Uri 'http://localhost:1323/app/notes'") */

	if *serve {
		core.Start() //blocking
	}

}

// TODO:?
// replace runServer()
// with:
// server = echo.New()
// core.RegisterRoutes(server)
