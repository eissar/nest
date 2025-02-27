package main

import (
	"flag"

	"github.com/eissar/nest/core"
)

// globals
var debug = false

func main() {

	help := flag.Bool("help", false, "print help information")
	serve := flag.Bool("serve", false, "run the utility server")
	//debug := flag.Bool("debug", true, "shows additional information in the console while running. (does nothing)")
	flag.Parse()

	if *help || flag.NFlag() == 0 {
		flag.PrintDefaults()
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
