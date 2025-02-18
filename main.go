package main

// dot source types
import (
	"errors"
	"log"
	_ "net"
	_ "sync"
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
	_ "github.com/eissar/nest/websocket-utils"

	_ "encoding/json"
	"flag"
	"fmt"
	"net/http"
	_ "net/http"

	_ "time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "golang.org/x/net/websocket"
)

// globals
var debug = false
var editor = "C:/Program Files/Neovim/bin/nvim.exe"

func runServer() {
	var err error
	server := echo.New()

	// TRAY ICON
	trayicon.Run(func() {
		core.Shutdown(server)
	})
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

	nestConfig := config.GetConfig()
	fmt.Printf(nestConfig.Host)

	// type dynamicTemplateHandlerOpts struct {
	// 	args  []string
	// 	first int
	// }
	// dynamicTemplateHandler := func(templateName string, populateFunc dynamicTemplatePopulateFunc, opts dynamicTemplateHandlerOpts) echo.HandlerFunc {

	// MIDDLEWARE
	// LOGGING
	noLog := []string{"/api/ping", "/template/open-tabs"}
	server.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Skipper: func(c echo.Context) bool {
			/* type Skipper func(c echo.Context) bool */
			for _, p := range noLog {
				if c.Path() == p {
					return true // (skip)
				}
			}
			return false
		},
		Format: "[LOG] [${time_rfc3339}] ${level} method=${method} path=${path}, Latency=${latency_human}\n",
	}))

	server.GET("/api/server/close", core.ServerShutdown)
	server.GET("/api/ping", core.Ping)

	// move somewhere else
	server.POST("/api/uploadTabs", core.UploadTabs)

	//server.GET("/api/broadcast/sse", broadcastHandler("getSong"))

	// Static routes;
	// route prefix, directory
	server.Static("css", "css")
	server.Static("js", "js")
	server.Static("img", "img")

	// Module philosophy:
	// UNDER NO CIRCUMSTANCES
	// should html or css be tightly coupled with
	// or packaged in a module (e.g., eagle_module) TODO:

	// ??? access routes in a module like:
	// server.GET("/eagleApp/*", eaglemodule.HandleModuleRoutes)
	// OR

	// routes for:
	// /eagle://item/<itemId>
	// /<itemId>
	eagle_module.RegisterRootRoutes(nestConfig, server)

	eagle_group := server.Group("/eagle")
	eagle_module.RegisterGroupRoutes(eagle_group)

	browser_group := server.Group("/browser")
	browser_module.RegisterGroupRoutes(browser_group)

	ytm_group := server.Group("/ytm")
	ytm_module.RegisterGroupRoutes(ytm_group)

	test_group := server.Group("/test")
	RegisterTestRoutes(test_group)

	template_group := server.Group("/template")
	RegisterTemplateRoutes(template_group)

	api_group := server.Group("/api")
	api.RegisterGroupRoutes(api_group)

	// special handler for user-facing static files
	// so file endings are not shown in the URI
	server.GET("/app/*", handlers.StaticAppHandler)

	if debug {
		//PrintSiteMap(server)
	}
	server.HideBanner = true

	err = server.Start(":1323")
	if err != nil {
		// CASE: server was closed by Server.Shutdown or Server.Close.
		if errors.Is(err, http.ErrServerClosed) {
			fmt.Println("[LOG] [SHUTDOWN] Shutting down gracefully...")
		} else {
			panic(err)
		}
	}
}

func trySearch() {
	e, err := eagle.New()
	if err != nil {
		log.Fatalf(err.Error())
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
}
