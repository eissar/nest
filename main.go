package main

// dot source types
import (
	"errors"
	_ "net"
	_ "sync"
	"web-dashboard/config"
	"web-dashboard/core"

	handlers "web-dashboard/handlers"

	browser_module "web-dashboard/modules/browser"
	eagle_module "web-dashboard/modules/eagle"
	ytm_module "web-dashboard/modules/ytm"

	render "web-dashboard/renderer-utils"
	_ "web-dashboard/websocket-utils"

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

// make api skip certain urls
// skipper = uri=/api/ping

// globals
var debug = false
var editor = "C:/Program Files/Neovim/bin/nvim.exe"

func runServer() {
	var err error
	server := echo.New()

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

	// server.GET("/", func(c echo.Context) error {
	// 	// fmt.Println(c.ParamNames())
	// 	if c.QueryParam("first") == "" {
	// 		c.QueryParams().Add("first", "5")
	// 	}
	// 	a := apiroutes.PopulateEnumerateWindows(c, "")
	// 	return c.Render(200, "windows.html", a)
	// })

	server.GET("/api/server/close", core.ServerShutdown)
	server.GET("/api/ping", core.Ping)

	// move somewhere else
	//server.GET("/api/windows", core.EnumWindows)
	server.GET("/api/numTabs", core.NumTabs)
	server.GET("/api/recentNotes", core.RecentNotes)
	server.POST("/api/edit", core.Edit)
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

	server.GET("/api/eagleOpen/:id", func(c echo.Context) error {
		id := c.Param("id")
		uri := fmt.Sprintf("eagle://item/%s", id)
		openURI(uri)
		return c.String(200, "OK")
	})

	// routes for:
	// /eagle://item/<itemId>
	// /<itemId>
	eagle_module.RegisterRootRoutes(server)

	eagle_group := server.Group("/eagle")
	eagle_module.RegisterRoutesFromGroup(eagle_group)

	browser_group := server.Group("/browser")
	browser_module.RegisterRoutesFromGroup(browser_group)

	ytm_group := server.Group("/ytm")
	ytm_module.RegisterRoutesFromGroup(ytm_group)

	test_group := server.Group("/test")
	RegisterTestRoutes(test_group)

	template_group := server.Group("/template")
	RegisterTemplateRoutes(template_group)

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

func main() {
	//#region parseFlags
	d := flag.Bool("debug", true, "shows additional information in the console while running.")
	flag.Parse()
	debug = *d
	//#endregion

	if debug {
		//pwsh.ExecPwshCmd("./powershell-utils/openUrl.ps1 -Uri 'http://localhost:1323/app/notes'")
	}
	runServer() //blocking
}
