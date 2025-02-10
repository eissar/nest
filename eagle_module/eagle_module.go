package eaglemodule

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"sync"
	_ "time"
	apiroutes "web-dashboard/api-routes"
	. "web-dashboard/types"

	_ "github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"golang.org/x/net/websocket"
	wsu "web-dashboard/websocket-utils"
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

//#REGION EXPERIMENTAL
// -------------------------------------------------------------------------------- //

// handle connections from our webserver <> webextension
// allows us to do things like query open browser tabs etc
// currently experimenting with straight querying the local db,
// but in that case, I have to implement that across browser archs.

// 	ws_connections = make(map[*websocket.Conn]bool)
// 	ws_mutex       sync.Mutex

// func initWebExtWS() *wsConfig {
// 	cfg := new(wsConfig)
//
// 	return nil
// }

func HandleWebExtWS(c echo.Context, cfg *wsu.WsConfig) error {
	// c.Response().Header().Set("Access-Control-Allow-Origin", "*")
	// this runs whenever a new connection request is made.

	// define a websocket handler.
	hndl := func(ws *websocket.Conn) {
		if len(cfg.Connections) >= 1 {
			fmt.Println("[WARN] You attempted to add another ws connection but maximum connection limit of 1 reached. ignoring this attempt. figure out how the background scripts work.")
			return
		}
		cfg.AddConnection(ws)
		defer func() {
			cfg.RemoveConnection(ws)
			ws.Close()
		}()

		//#region read
		for {
			msg := ""
			err := websocket.Message.Receive(ws, &msg)
			if err != nil {
				if err.Error() == "EOF" {
					// fmt.Println("[LOG] <HandleWS> Client disconnected (EOF)")
					// cfg.RemoveConnection(ws)
					return // stop serving.
				} else if errors.Is(err, io.EOF) {
					// fmt.Println("[LOG] <HandleWS> Client disconnected (EOF)")
					// cfg.RemoveConnection(ws)
					return // stop serving.
				} else {
					// fmt.Println("[ERROR] <HandleWS> Unexpected error receiving message:", err)
					// cfg.RemoveConnection(ws)
					return // stop serving.
				}
			}

			fmt.Println("[LOG] <HandleWebExtWS> Received: ", msg)
			// if the message has an id field, process it.

			fmt.Println("client says", msg)

			// Broadcast the received message to all clients
			// broadcastMessage := fmt.Sprintf("Client says: %s", msg)
			// Broadcast(broadcastMessage)
		}
		//#endregion read
	}
	// handle the websocket
	websocket.Handler(hndl).ServeHTTP(
		c.Response(),
		c.Request(),
	)

	// for conn := range ws.Connections {
	// 	ws.
	// }

	return nil
}

// func GetResponseWS(c echo.Context, ws *wsu.WsConfig) error {
// 	a := pendingRequests.New()
// 	a.Id
//
// 	return nil
// }

// -------------------------------------------------------------------------------- //
//#ENDREGION EXPERIMENTAL

func RegisterRoutesFromGroup(g *echo.Group) {
	g.GET("/:id", func(c echo.Context) error {
		id := c.Param("id")
		out := fmt.Sprintf("<p>the id is : %s!</p>", id)
		return c.String(200, out)
	})

	//#region sse
	g.GET("/sse", HandleSSE)
	g.GET("/test-broadcast-to-sse-clients", func(c echo.Context) error {
		apiroutes.SendSSETargeted(`{"data": "test","event": "getTabs"}`, clients, &clientsMu)
		//apiroutes.SendSSETargeted("data: test\nevent: message", clients, &clientsMu)
		return c.String(200, "OK")
	})
	//#endregion sse

	//#region ws
	websocketCfg := wsu.WsConfig{
		Connections: make(map[*websocket.Conn]bool),
		Mutex:       sync.Mutex{},
	}
	g.GET("/ws", func(c echo.Context) error {
		err := HandleWebExtWS(c, &websocketCfg)
		if err != nil {
			fmt.Println("[ERROR]", err.Error())
			c.String(400, err.Error())
		}
		return c.String(200, "OK")
	})
	g.GET("/wsbroad", func(c echo.Context) error {
		websocketCfg.Broadcast(`{"id":"1","command":"tabsCount"}`)
		// add Pending WSRequest with callback?
		return c.String(200, "OK")
	})
	//#endregion ws
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
