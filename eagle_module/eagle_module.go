package eaglemodule

import (
	"fmt"

	"github.com/labstack/echo/v4"
)

func HandleModuleRoutes(c echo.Context) error {
	requestPath := c.Param("*") // Get the requested path after "/app/"
	fmt.Println("requestPath:", requestPath)
	return c.String(200, "OK")
}

func RegisterRoutesFromGroup(g *echo.Group) {
	g.GET("/:id", func(c echo.Context) error {
		id := c.Param("id")

		out := fmt.Sprintf("<p>the id is : %s!</p>", id)
		return c.String(200, out)
	})

}

/*
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

*/
