package main

// dot source types
import (
	"errors"
	_ "net"
	"os/exec"
	"path/filepath"
	"runtime"
	_ "sync"
	apiroutes "web-dashboard/api-routes"
	browser_module "web-dashboard/browser_module"
	databaseutils "web-dashboard/database-utils"
	eagle_module "web-dashboard/eagle_module"
	pwsh "web-dashboard/powershell-utils"
	_ "web-dashboard/websocket-utils"
	ytm_module "web-dashboard/ytm_module"

	_ "encoding/json"
	"flag"
	"fmt"
	"net/http"
	_ "net/http"

	"strings"
	_ "time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "golang.org/x/net/websocket"
)

/* __MACRO__
```lua
!wt.exe -d "$env:CLOUD_DIR\Code\go\web-dashboard" pwsh -c ./build.ps1
```
:so ./build.lua

Explanation:
open new terminal instance with the server running.
*/

/*
// TODO:
			[X] - move types to types.go
			[X] - add HTMX
			[X] - add middleware to serve url without file suffix.
			[X] - Recent notes
			[X] - move time.now calls to middleware (custom)
			q5s: enumerate-Windows.ps1
			[ ] - Try reflection for template functions?
			[ ] - add action parameter to recent notes
			[X] - Create build.lua
			[X] - add sse listener to eagle-plugin
				[X] - for eagle.tabs.query({})
				on BroadcastSSETargeted({event:"getTabs"...})
					-> eagle.tabs.query({})
					-> post("api/uploadTabs")
					-> fmt.PrintLn tabs
*/

// make api skip certain urls
// skipper = uri=/api/ping

// globals
var debug = false
var editor = "C:/Program Files/Neovim/bin/nvim.exe"

/*
Inputs: path to .ps1 script
Outputs: array-contained json data.
*/

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
	server.Renderer = &Template{
		templates: mustImportTemplates(),
	}

	// HANDLERS
	// handler closures for satisfying echo.HandlerFunc signature so this can be pretty
	// static templates
	staticTemplateHandler := func(templateName string) echo.HandlerFunc {
		return func(c echo.Context) error {
			return c.Render(200, templateName, nil)
		}
	}
	// special handler for user-facing static files
	// so file endings are not shown in the URI
	staticAppHandler := func(c echo.Context) error {
		// Serve static files with fallback for /app/index
		requestPath := c.Param("*") // Get the requested path after "/app/"
		isFullPath := strings.HasSuffix(requestPath, ".html")
		if !isFullPath {
			requestPath = requestPath + ".html"
		}
		filePath := filepath.Join("html", requestPath)
		return c.File(filePath)
	}
	// broadcastHandler := func(broadcastName string) echo.HandlerFunc {
	// 	return func(c echo.Context) error {
	// 		return BroadcastEvent(c, broadcastName)
	// 	}
	// }
	// dynamicTemplatePopulateFunc defines a function to populate a template
	type dynamicTemplatePopulateFunc func(c echo.Context, templateName string) interface{}

	// type dynamicTemplateHandlerOpts struct {
	// 	args  []string
	// 	first int
	// }
	// dynamicTemplateHandler := func(templateName string, populateFunc dynamicTemplatePopulateFunc, opts dynamicTemplateHandlerOpts) echo.HandlerFunc {

	// closure generator
	// returns echo.HandlerFunc
	// uses populateFunc to populate template with template name (incl. ending)
	// opts are { args: []string{} }
	dynamicTemplateHandler := func(templateName string, populateFunc dynamicTemplatePopulateFunc) echo.HandlerFunc {
		// dynamicTemplatePopulateFunc
		return func(c echo.Context) error {
			// to set default parameters, update them in the populateFunc.
			return c.Render(200, templateName, populateFunc(c, templateName))
		}
	}
	type pwshTemplateType string
	const (
		pwshScript  pwshTemplateType = "pwshScript"
		pwshCommand pwshTemplateType = "pwshCommand"
	)
	pwshTemplateHandler := func(templateName string, typ pwshTemplateType, p string) echo.HandlerFunc {
		if typ != pwshScript {
			panic("yeah")
		}
		return func(c echo.Context) error {
			return c.Render(200, templateName, pwsh.RunPwshCmd(p))
		}
	}

	// MIDDLEWARE
	// LOGGING
	server.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Skipper: func(c echo.Context) bool {
			/* type Skipper func(c echo.Context) bool */
			if c.Path() == "/api/ping" {
				return true // (skip)
			}
			if c.Path() == "/template/open-tabs" {
				return true // (skip)
			}
			return false
		},
		Format: "[LOG] [${time_rfc3339}] ${level} method=${method} path=${path}, Latency=${latency_human}\n",
	}))

	//#region TEST
	/*
		try the renderer with recentNotes
		test := func(c echo.Context) error {
			a := apiroutes.GetNotesNamesDates()
			fmt.Println(a[0])
			return c.String(200, "OK")

		}
	*/
	// server.GET("/test", test)
	server.GET("/test", dynamicTemplateHandler("notes-struct.html", apiroutes.PopulateGetNotesDetail))

	server.GET("/test1", func(c echo.Context) error {
		db, err := databaseutils.GetDatabase()
		if err != nil {
			fmt.Println("error calling <GetDatabase>", err)
			return c.String(400, "NOT OK")
		}
		fmt.Println(db.Data)
		return c.String(200, "OK")
	})
	//#endregion TEST

	server.GET("/", func(c echo.Context) error {
		// fmt.Println(c.ParamNames())
		if c.QueryParam("first") == "" {
			c.QueryParams().Add("first", "5")
		}
		a := apiroutes.PopulateEnumerateWindows(c, "")
		return c.Render(200, "windows.html", a)
	})
	// WEBSOCKET
	//server.GET("/ws", websockettest.HandleWS)

	// SERVER SIDE EVENTS
	//server.GET("/sse", HandleSSE)
	// server.GET("/eagle/sse", HandleSSE) in eaglemodule

	server.GET("/template/notes-struct", dynamicTemplateHandler("notes-struct.html", apiroutes.PopulateGetNotesDetail))
	server.GET("/template/windows", dynamicTemplateHandler("windows.html", apiroutes.PopulateEnumerateWindows))

	server.GET("/template/recent-notes", pwshTemplateHandler("recent-notes.html", pwshScript, "./powershell-utils/recentNotes.ps1"))
	server.GET("/template/key-value", pwshTemplateHandler("key-value.templ", pwshScript, "./powershell-utils/mock_nvim.ps1"))
	server.GET("/template/open-tabs-count", pwshTemplateHandler("open-tabs-count.templ", pwshScript, "./powershell-utils/waterfoxTabs.ps1"))
	//server.GET("/template/open-tabs", dynamicTemplateHandler("open-tabs.templ", apiroutes.PopulateOpenTabs))

	server.GET("/template/recent-eagle-items", pwshTemplateHandler("recent-eagle-items.templ", pwshScript, "./powershell-utils/recentEagleItems.ps1"))

	server.GET("/template/sse-browser-tabs", staticTemplateHandler("sse-browser-tabs.templ"))

	server.GET("/template/browser-tabs", staticTemplateHandler("browser-tabs.templ"))

	server.GET("/template/recent-notes_layout", staticTemplateHandler("recent-notes.layout.html"))
	server.GET("/template/timeline_layout", staticTemplateHandler("timeline.layout.html"))
	server.GET("/template/now-playing", staticTemplateHandler("ws-now-playing.ytm.templ")) // ./templates/ws-now-playing.ytm.templ
	server.GET("/template/ping", staticTemplateHandler("ping.templ"))

	server.GET("/api/server/close", apiroutes.ServerShutdown)
	server.GET("/api/windows", apiroutes.EnumWindows)
	server.GET("/api/numTabs", apiroutes.NumTabs)
	server.GET("/api/recentNotes", apiroutes.RecentNotes)
	server.GET("/api/recentEagleItems", apiroutes.RecentEagleItems)
	server.POST("/api/edit", apiroutes.Edit)
	server.GET("/api/ping", apiroutes.Ping)
	server.POST("/api/uploadTabs", apiroutes.UploadTabs)
	// activates a browser tab if it exists, creates a new tab if it does not.
	// cache is created from browser history
	// api.BrowserTabActivateOrOpen

	// api.GetBrowserHistory

	// about:profiles
	// ["X:\Dropbox\Code\Projects\render-image-blazingly\db"]

	//server.GET("/api/broadcast/sse", broadcastHandler("getSong"))

	openURI := func(uri string) error {
		var cmd *exec.Cmd

		fmt.Println("[LOG] <openUri> opening...", uri)
		switch runtime.GOOS {
		case "windows":
			cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", uri)
		case "darwin":
			cmd = exec.Command("open", uri)
		default: // Linux and other Unix-like systems
			cmd = exec.Command("xdg-open", uri)
		}

		return cmd.Run()
	}

	server.GET("/api/eagleOpen/:id", func(c echo.Context) error {
		id := c.Param("id")
		uri := fmt.Sprintf("eagle://item/%s", id)
		openURI(uri)
		return c.String(200, "OK")
	})

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

	eagle_group := server.Group("/eagleApp")
	eagle_module.RegisterRoutesFromGroup(eagle_group)

	browser_group := server.Group("/browser")
	browser_module.RegisterRoutesFromGroup(browser_group)

	ytm_group := server.Group("/ytm")
	ytm_module.RegisterRoutesFromGroup(ytm_group)

	// special handler for user-facing static files
	// so file endings are not shown in the URI
	server.GET("/app/*", staticAppHandler)

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
