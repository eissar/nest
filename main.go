package main

import (
	"flag"
	"github.com/eissar/nest/core"
)

// globals
var debug = false

func main() {
	//#region parseFlags
	d := flag.Bool("debug", true, "shows additional information in the console while running. (does nothing)")
	flag.Parse()
	debug = *d
	//#endregion
	if debug {
		/* pwsh.ExecPwshCmd("./powershell-utils/openUrl.ps1 -Uri 'http://localhost:1323/app/notes'") */
	}

	// trySearch()
	core.Start() //blocking

	// TODO:?
	// replace runServer()
	// with:
	// server = echo.New()
	// core.RegisterRoutes(server)
}
