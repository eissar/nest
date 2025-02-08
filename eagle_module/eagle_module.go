package eaglemodule

import (
	"fmt"
	"net/http"
	"sync"
	apiroutes "web-dashboard/api-routes"
	. "web-dashboard/types"

	"github.com/labstack/echo/v4"
)

func HandleModuleRoutes(c echo.Context) error {
	requestPath := c.Param("*") // Get the requested path after "/app/"
	fmt.Println("requestPath:", requestPath)
	return c.String(200, "OK")
}

var (
	clientsMu sync.RWMutex
	clients   = make(map[*Client]bool)
)

func HandleSSE(c echo.Context) error {
	c.Response().Header().Set("Access-Control-Allow-Origin", "*")
	c.Response().Header().Set("Content-Type", "text/event-stream")
	c.Response().Header().Set("Cache-Control", "no-cache")
	c.Response().Header().Set("Connection", "keep-alive")

	flusher, ok := c.Response().Writer.(http.Flusher)
	if !ok {
		return fmt.Errorf("streaming unsupported")
	}

	client := &Client{
		Conn: c.Response(),
		Ch:   make(chan string),
	}

	clientsMu.Lock()
	clients[client] = true
	clientsMu.Unlock()

	defer func() {
		clientsMu.Lock()
		delete(clients, client)
		clientsMu.Unlock()
		close(client.Ch)
		fmt.Println("[LOG] <HandleSSE> Client disconnected")
	}()

	go func() {
		for msg := range client.Ch {
			fmt.Fprintf(client.Conn, "data: %s\n\n", msg)
			flusher.Flush()
		}
	}()

	// Keep the connection open.  No need to read from the client.
	<-c.Request().Context().Done() // Block until client disconnects
	fmt.Println("[LOG] <eagle_module.HandleSSE> Client disconnected")
	return nil
}

func RegisterRoutesFromGroup(g *echo.Group) {
	g.GET("/:id", func(c echo.Context) error {
		id := c.Param("id")
		out := fmt.Sprintf("<p>the id is : %s!</p>", id)
		return c.String(200, out)
	})

	g.GET("/sse", HandleSSE)
	g.GET("/test-broadcast-to-sse-clients", func(c echo.Context) error {
		apiroutes.SendSSETargeted(`{"data": "test","event": "getTabs"}`, clients, &clientsMu)
		//apiroutes.SendSSETargeted("data: test\nevent: message", clients, &clientsMu)
		return c.String(200, "OK")
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
