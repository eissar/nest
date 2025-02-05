package main

// dot source types
import (
	"context"
	"errors"
	_ "net"
	"path/filepath"
	_ "sync"
	"time"
	apiroutes "web-dashboard/api-routes"
	pwsh "web-dashboard/powershell-utils"

	_ "encoding/json"
	"flag"
	"fmt"
	"html/template"
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
//	TODO:
//	[X] - move types to types.go
//	[X] - add HTMX
//	[X] - add middleware to serve url without file suffix.
//  [X] - Recent notes
//	[ ] - move time.now calls to middleware (custom)
//  [ ] - add action to recent notes
// 	q5s: enumerate-Windows.ps1
//  [ ] - Try reflection for template functions.
//  [X] - Create build.lua
*/

var debug = false
var editor = "C:/Program Files/Neovim/bin/nvim.exe"
var lastSong = "NULL SONG DATA"

/*
//	Inputs: path to .ps1 script
//	Outputs: array-contained json data.
*/
func mustImportTemplates() *template.Template {
	templ, err := template.ParseGlob("templates/*") // Parses all .html files in the templates directory
	if err != nil {
		panic(err)
	}
	return templ
}
func shutdownServerWithTimeout(timeoutCtx context.Context, server *echo.Echo) {
	err := server.Shutdown(timeoutCtx)
	if err != nil {
		panic(err)
	}
}

func runServer() {
	var err error
	templ := mustImportTemplates()
	server := echo.New()

	server.GET("/", func(c echo.Context) error {
		a := (apiroutes.GetEnumerateWindows()[0])

		err = templ.ExecuteTemplate(c.Response().Writer, "test.html", a)
		if err != nil {
			fmt.Println(err)
			return c.String(http.StatusInternalServerError, "Error rendering template, check logs")
		}

		return nil
	})
	server.GET("/template/template-test", func(c echo.Context) error {
		a := apiroutes.GetEnumerateWindows()[0]
		err = templ.ExecuteTemplate(c.Response().Writer, "window.html", a)
		return nil
	})
	server.GET("/template/windows", func(c echo.Context) error {
		a := apiroutes.GetEnumerateWindows()
		err = templ.ExecuteTemplate(c.Response().Writer, "windows.html", a)
		return nil
	})
	server.GET("/template/recent-notes", func(c echo.Context) error {
		a := pwsh.RunPwshCmd("./recentNotes.ps1")
		err = templ.ExecuteTemplate(c.Response().Writer, "recent-notes.html", a)
		return nil
	})
	server.GET("/template/recent-notes_layout", func(c echo.Context) error {
		err = templ.ExecuteTemplate(c.Response().Writer, "recent-notes.layout.html", nil)
		return nil
	})
	server.GET("/template/timeline_layout", func(c echo.Context) error {
		err = templ.ExecuteTemplate(c.Response().Writer, "timeline.layout.html", nil)
		return nil
	})
	server.GET("/template/key-value", func(c echo.Context) error {
		a := pwsh.RunPwshCmd("./mock_nvim.ps1")
		err = templ.ExecuteTemplate(c.Response().Writer, "key-value.templ", a)
		if err != nil {
			fmt.Println(err)
		}
		return nil
	})
	server.GET("/template/ping", func(c echo.Context) error {
		err = templ.ExecuteTemplate(c.Response().Writer, "ping.html", nil)
		return nil
	})
	server.GET("/template/now-playing", func(c echo.Context) error {
		err = templ.ExecuteTemplate(c.Response().Writer, "now-playing.static.html", nil)
		return nil
	})
	server.GET("/template/open-tabs", func(c echo.Context) error {
		a := pwsh.RunPwshCmd("./waterfoxTabs.ps1")
		err = templ.ExecuteTemplate(c.Response().Writer, "open-tabs.static.html", a)
		return nil
	})

	server.GET("/api/server/close", func(c echo.Context) error {
		//use go routine with timeout to allow time for response.
		var err error

		timeout := 10 * time.Second
		timeoutCtx, shutdownRelease := context.WithTimeout(context.Background(), timeout)
		defer shutdownRelease()
		// go shutdownServerWithTimeout(timeoutCtx, server)
		// try doing the shutdown inline without goroutine.

		go func() {
			err = server.Shutdown(timeoutCtx)
		}()
		if err != nil {
			fmt.Println("err while graceful shutdown:", err)
		}

		return c.String(200, "shutdown cmd successful.")
	})
	server.GET("/api/windows", apiroutes.EnumWindows)
	server.GET("/api/numTabs", apiroutes.NumTabs)
	server.GET("/api/recentNotes", apiroutes.RecentNotes)
	server.GET("/api/recentEagleItems", apiroutes.RecentEagleItems)
	server.POST("/api/edit", apiroutes.Edit)
	server.GET("/api/ping", apiroutes.Ping)

	server.GET("/api/broadcast/yt-music", func(c echo.Context) error {
		return BroadcastEvent(c, "ytMusicElement")
	})
	server.GET("/api/broadcast/sse", func(c echo.Context) error {
		return BroadcastEvent(c, "getSong")
	})

	// WEBSOCKET
	server.GET("/ws", Hello)
	// SERVER SIDE EVENTS
	server.GET("/sse", HandleSSE)

	// route prefix, directory
	server.Static("css", "css")
	server.Static("js", "js")
	server.Static("img", "img")
	server.GET("/app/*", func(c echo.Context) error {
		// Serve static files manually with fallback for /app/index
		requestPath := c.Param("*") // Get the requested path after "/app/"
		isFullPath := strings.HasSuffix(requestPath, ".html")
		if !isFullPath {
			requestPath = requestPath + ".html"
		}
		filePath := filepath.Join("html", requestPath)
		return c.File(filePath)
	})
	//server.GET("/template/*", func(c echo.Context) error {
	//	requestPath := c.Param("*") // Get the requested path after "/app/"
	//	filePath := filepath.Join("html", requestPath)
	//	return nil
	//})

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
		Format: "[LOG] [${time_rfc3339}] ${level} path=${path}, Latency=${latency_human}\n",
	}))
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
	d := flag.Bool("debug", true, "shows additional information in the console while running.")
	flag.Parse()
	debug = *d

	if debug {
		pwsh.ExecPwshCmd("./openUrl.ps1 -Uri 'http://localhost:1323'")
	}
	runServer()
	// when starting the server, send SSE to all yt-music clients if music is playing.
}
