package handlers

import (
	"path/filepath"
	"strings"

	"github.com/eissar/nest/plugins/pwsh"
	render "github.com/eissar/nest/render"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// handlers for satisfying echo.HandlerFunc
// and middlewareFunc

type dynamicPopulateFunc = render.DynamicTemplatePopulateFunc

// handlers for satisfying echo.HandlerFunc
// and closures for generating the same.

// static templates
func StaticTemplateHandler(templateName string) echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.Render(200, templateName, nil)
	}
}

// StaticAppHandler :=
// special handler for user-facing static files
// so file endings are not shown in the URI
func StaticAppHandler(c echo.Context) error {
	// Serve static files with fallback for /app/index
	requestPath := c.Param("*") // Get the requested path after "/app/"
	isFullPath := strings.HasSuffix(requestPath, ".html")
	if !isFullPath {
		requestPath = requestPath + ".html"
	}
	filePath := filepath.Join("html", requestPath)
	return c.File(filePath)
}

// type dynamicTemplateHandlerOpts struct {
// 	args  []string
// 	first int
// }
// dynamicTemplateHandler := func(templateName string, populateFunc dynamicTemplatePopulateFunc, opts dynamicTemplateHandlerOpts) echo.HandlerFunc {

// closure generator
// returns echo.HandlerFunc
// uses populateFunc to populate template with template name (incl. ending)
// opts are { args: []string{} }
func DynamicTemplateHandler(templateName string, populateFunc dynamicPopulateFunc) echo.HandlerFunc {
	// dynamicTemplatePopulateFunc
	return func(c echo.Context) error {
		// to set default parameters, update them in the populateFunc.
		return c.Render(200, templateName, populateFunc(c, templateName))
	}
}

func PwshTemplateHandler(templateName string, p string) echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.Render(200, templateName, pwsh.RunScript(p))
	}
}

// param excl []string excluded paths (paths we don't log)
func LoggerMiddleware(excl []string) echo.MiddlewareFunc {
	return middleware.LoggerWithConfig(middleware.LoggerConfig{
		Skipper: func(c echo.Context) bool {
			/* type Skipper func(c echo.Context) bool */
			for _, p := range excl {
				if c.Path() == p {
					return true // (skip)
				}
			}
			return false
		},
		Format: "[LOG] [${time_rfc3339}] ${level} method=${method} path=${path}, Latency=${latency_human}\n",
	})
}

// old
// handler closures for satisfying echo.HandlerFunc signature so this can be pretty
// static templates
/*
	staticTemplateHandler := func(templateName string) echo.HandlerFunc {
		return func(c echo.Context) error {
			return c.Render(200, templateName, nil)
		}
	}
*/
//staticTemplateHandler :=
// special handler for user-facing static files
// so file endings are not shown in the URI
// staticAppHandler := func(c echo.Context) error {
// 	// Serve static files with fallback for /app/index
// 	requestPath := c.Param("*") // Get the requested path after "/app/"
// 	isFullPath := strings.HasSuffix(requestPath, ".html")
// 	if !isFullPath {
// 		requestPath = requestPath + ".html"
// 	}
// 	filePath := filepath.Join("html", requestPath)
// 	return c.File(filePath)
// }
// broadcastHandler := func(broadcastName string) echo.HandlerFunc {
// 	return func(c echo.Context) error {
// 		return BroadcastEvent(c, broadcastName)
// 	}
// }
