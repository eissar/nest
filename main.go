package main

// dot source types
import (
	"errors"
	"log"
	"time"

	"github.com/eissar/nest/config"
	"github.com/eissar/nest/core"
	trayicon "github.com/eissar/nest/core/tray-icon"
	"github.com/eissar/nest/eagle"
	"github.com/eissar/nest/eagle/api"

	handlers "github.com/eissar/nest/handlers"

	browser_module "github.com/eissar/nest/plugins/browser"
	eagle_module "github.com/eissar/nest/plugins/eagle"
	"github.com/eissar/nest/plugins/search"
	ytm_module "github.com/eissar/nest/plugins/ytm"

	"github.com/eissar/nest/render"

	"flag"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

// globals
var debug = false
var editor = "C:/Program Files/Neovim/bin/nvim.exe"

func runServer() {
	var err error

	nestConfig := config.GetConfig()

	server := echo.New()

	// TRAY ICON
	trayicon.Run(server)
	defer trayicon.Quit()

	// NOTE: Template rules:
	// 1. ending in .html:  static template.
	// 2. ending in .templ:  dynamic template
	// 3. prefix ws, sse, ending in .templ:  template which retrieves data...
	// using websockets or server-side events respectively.
	// 4. no prefix, ending in .templ: template which retrieves data...
	// dynamically using htmx.
	server.Renderer = &render.Template{
		Templates: render.MustImportTemplates(),
	}

	fmt.Printf("%s", nestConfig.Host)

	// MIDDLEWARE LOGGING
	excludedPaths := []string{"/api/ping", "/template/open-tabs"}
	server.Use(handlers.LoggerMiddleware(excludedPaths))

	// SCOPED ROUTES
	eagle_group := server.Group("/eagle")
	eagle_module.RegisterGroupRoutes(eagle_group)

	eagleapi_group := server.Group("/api")
	api.RegisterGroupRoutes(eagleapi_group)

	browser_group := server.Group("/browser")
	browser_module.RegisterGroupRoutes(browser_group)

	ytm_group := server.Group("/ytm")
	ytm_module.RegisterGroupRoutes(ytm_group)

	test_group := server.Group("/test")
	RegisterTestRoutes(test_group)

	template_group := server.Group("/template")
	RegisterTemplateRoutes(template_group)

	// ROOT ROUTES
	eagle_module.RegisterRootRoutes(nestConfig, server)
	core.RegisterRootRoutes(server)
	api.RegisterRootRoutes(server)

	// STATIC ROUTES (route prefix, directory)
	server.Static("css", "css")
	server.Static("js", "js")
	server.Static("img", "img")

	// special handler for user-facing static files
	// so file endings are not shown in the URI
	server.GET("/app/*", handlers.StaticAppHandler)

	if debug {
		//PrintSiteMap(server)
	}
	server.HideBanner = true

	err = server.Start(":1323")
	if err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			// shutdown was requested
			fmt.Println("[LOG] [SHUTDOWN] Shutting down gracefully...")
		} else {
			// crash
			panic(err)
		}
	}
}

func trySearch() {
	e, err := eagle.New()
	if err != nil {
		log.Fatalf("%s", err.Error())
	}
	s := search.New(e)
	defer s.Index.Close()

	//go search.Index(e, s.Index)
	//search.ForceReIndex(e, s.Index)
	search.ForceReIndexStreaming(e, s.Index)
	return

	start := time.Now()
	s.Query("vallejo")
	fmt.Print("search took: ", time.Since(start))
}

func main() {
	//#region parseFlags
	d := flag.Bool("debug", true, "shows additional information in the console while running.")
	flag.Parse()
	debug = *d
	//#endregion

	if debug {
		// pwsh.ExecPwshCmd("./powershell-utils/openUrl.ps1 -Uri 'http://localhost:1323/app/notes'")
	}

	// trySearch()
	runServer() //blocking

	// TODO:?
	// replace runServer()
	// with:
	// server = echo.New()
	// core.RegisterRoutes(server)
}
