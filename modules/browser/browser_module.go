package browsermodule

import (
	"fmt"
	"sync"
	"time"
	_ "time"
	//browsermodule "github.com/eissar/nest/browser_module"
	//. "github.com/eissar/nest/types"

	wsu "github.com/eissar/nest/websocket-utils"

	_ "github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/net/websocket"
)

func HandleModuleRoutes(c echo.Context) error {
	requestPath := c.Param("*") // Get the requested path after "/app/"
	fmt.Println("requestPath:", requestPath)
	return c.String(200, "OK")
}

var (
	DefaultTimeout = 3 * time.Second
)

// handle connections from our webserver <> webextension
// allows us to do things like query open browser tabs etc
// currently experimenting with straight querying the local db,
// but in that case, I have to implement that across browser archs.

var PendingRequests wsu.PendingWSRequests
var PendingRequestsV2 wsu.PendingWSRequestsV2

func RegisterGroupRoutes(g *echo.Group) {
	trimmer := time.NewTicker(5 * time.Minute)
	defer trimmer.Stop()
	go func() {
		for range trimmer.C {
			PendingRequests.Trim()
		}
	}()

	//#region sse
	g.GET("/sse", HandleSSE)
	// g.GET("/test-broadcast-to-sse-clients", func(c echo.Context) error {
	// 	core.SendSSETargeted(`{"data": "test","event": "getTabs"}`, clients, &clientsMu)
	// 	//apiroutes.SendSSETargeted("data: test\nevent: message", clients, &clientsMu)
	// 	return c.String(200, "OK")
	// })
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
	api := g.Group("/api")
	api.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		Skipper:      middleware.DefaultSkipper,
		Timeout:      8 * time.Second,
		ErrorMessage: "api middleware: internal timeout (8s) has been reached. you should not be seeing this.",

		OnTimeoutRouteErrorHandler: func(err error, c echo.Context) {
			fmt.Println("[ERROR] Timeout reached in route:", c.Request().URL, "params:", c.QueryParams())
		},
	}))

	api.GET("/getTabs", HandleFilteredGetTabs(&websocketCfg))
	// depends on accessbile webExt
	g.GET("/broadcast/getTabs", HandleGetTabs(&websocketCfg))

	g.GET("/broadcast/tabsCount", func(c echo.Context) error {
		websocketCfg.Broadcast(`{"id":"1","command":"tabsCount"}`)
		// add Pending WSRequest with callback?
		return c.String(200, "OK")
	})
	/* old
	g.GET("/broadcast/getTabs", func(c echo.Context) error {
		req := PendingRequests.New()
		bc := fmt.Sprintf(`{"id":"%s","command":"getTabs"}`, req.Id)
		websocketCfg.Broadcast(bc)
		// add Pending WSRequest with callback?
		return c.String(200, "OK")
	})
	*/

	/*
		g.GET("/broadcast/getTabs", func(c echo.Context) error {
			// test if the webext is connected.
			if len(websocketCfg.Connections) < 1 {
				resp := fmt.Sprintf("error=there are no connections to a browser extension. "+
					"first run the browser extension to enable retrieving browser data "+
					"at request to :%s", c.Path())
				return c.String(400, resp)
			}

			// create a timeout
			timer := time.NewTimer(DefaultTimeout)
			defer timer.Stop()

			// init a new request
			req := PendingRequestsV2.New()
			//fmt.Println("added request:", req.Id)

			// read: handleWebExtWS will match id
			// in PendingRequestsV2 when a
			// message is recieved and sends data to
			// the channel on successful match.

			// broadcast a request for tabs
			// clients that recieve this will try to
			// respond with a matching id
			bc := fmt.Sprintf(`{"id":"%s","command":"getTabs"}`, req.Id)
			websocketCfg.Broadcast(bc)

			select {
			//case _, ok := <-timer.C:
			case resp := <-req.Chan:
				return c.JSON(200, resp)
			case _ = <-timer.C:
				// timeout has elapsed
				resp := fmt.Sprintf("error: timeout reached (%s), at request to :%s", DefaultTimeout, c.Path())
				return c.String(400, resp)
			}
			return c.String(http.StatusInternalServerError, "FALSE")
		})
	*/

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
