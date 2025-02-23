package core

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/eissar/nest/config"
	trayicon "github.com/eissar/nest/core/tray-icon"
	"github.com/eissar/nest/eagle/api"
	"github.com/eissar/nest/handlers"
	nest "github.com/eissar/nest/plugins/eagle"
	"github.com/eissar/nest/render"
	"github.com/eissar/nest/templates"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func Start() {
	var err error

	nestConfig := config.GetConfig()

	if isPortOccupied(nestConfig.Nest.Port) {
		log.Fatalf("error starting server: port %d occupied. is the server already running?", nestConfig.Nest.Port)
	}

	server := echo.New()

	// TRAY ICON
	trayicon.Run(server, func() {
		Shutdown(server) // onExit trayicon function
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

	fmt.Printf("%s", nestConfig.Host)

	// MIDDLEWARE LOGGING
	excludedPaths := []string{"/api/ping", "/template/open-tabs"}
	server.Use(handlers.LoggerMiddleware(excludedPaths))

	server.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"app://obsidian.md"},
	}))

	// SCOPED ROUTES
	eagle_group := server.Group("/eagle")
	nest.RegisterGroupRoutes(eagle_group)

	eagleapi_group := server.Group("/api")
	api.RegisterGroupRoutes(eagleapi_group)

	template_group := server.Group("/template")
	templates.RegisterTemplateRoutes(template_group)

	// ROOT ROUTES
	nest.RegisterRootRoutes(nestConfig, server)
	api.RegisterRootRoutes(server)
	RegisterRootRoutes(server)

	// STATIC ROUTES (route prefix, directory)
	server.Static("css", "./assets/css")
	server.Static("js", "./assets/js")
	server.Static("img", "./assets/img")

	// special handler for user-facing static files
	// so file endings don't have to be shown in the URI
	server.GET("/app/*", handlers.StaticAppHandler)

	server.HideBanner = true

	err = server.Start(fmt.Sprintf(":%d", nestConfig.Nest.Port))
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

func PrintSiteMap(server *echo.Echo) {
	fmt.Println("server available routes:")
	for _, x := range server.Routes() {
		fmt.Println(x.Name, x.Path)
	}
}

func Ping(c echo.Context) error {
	return c.String(http.StatusOK, "OK")
}

func isPortOccupied(port int) bool {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return true // Port is likely occupied
	}
	defer listener.Close()
	return false // Port is free
}

// func PopulateQueryFrontmatter() {}
// func PopulateGetRecentNotes()   {}

func ServerShutdown(c echo.Context) error {
	var err error

	//use go routine with timeout to allow time for response.
	timeout := 10 * time.Second
	timeoutCtx, shutdownRelease := context.WithTimeout(context.Background(), timeout)
	defer shutdownRelease()

	go func() {
		err = c.Echo().Server.Shutdown(timeoutCtx)
	}()
	if err != nil {
		fmt.Println("err while graceful shutdown:", err)
	}

	return c.String(200, "shutdown cmd successful.")
}
func Shutdown(s *echo.Echo) error {
	var err error
	//use go routine with timeout to allow time for response.
	timeout := 10 * time.Second
	timeoutCtx, shutdownRelease := context.WithTimeout(context.Background(), timeout)
	defer shutdownRelease()

	go func() {
		err = s.Shutdown(timeoutCtx)
	}()
	if err != nil {
		fmt.Println("err while graceful shutdown:", err)
	}
	return nil
}

// registers routes on the server root (/)
func RegisterRootRoutes(server *echo.Echo) {
	//server.GET("/eagle\\://item/:itemId", ServeThumbnailHandler(&n))
	server.GET("/api/server/close", ServerShutdown)
	server.GET("/api/ping", Ping)
}
