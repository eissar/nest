package browsermodule

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/eissar/nest/modules/pwsh"
	. "github.com/eissar/nest/types"
	wsu "github.com/eissar/nest/websocket-utils"

	"github.com/gobwas/glob"
	"github.com/labstack/echo/v4"
	"golang.org/x/net/websocket"
)

// closure
// depends on accessible webExt
func HandleGetTabs(w *wsu.WsConfig) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set("Access-Control-Allow-Origin", "*")
		// test if the webext is connected.
		if len(w.Connections) < 1 {
			resp := fmt.Sprintf("error=there are no connections to a browser extension. "+
				"first run the browser extension to enable retrieving browser data "+
				"at request to :%s", c.Path())
			return c.String(400, resp)
		}

		// create a timeout
		timer := time.NewTimer(3 * time.Second)
		//defer timer.Stop()

		//wsu.PendingWSRequestsV2.Mutex.Lock()
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
		w.Broadcast(bc)

		select {
		//case _, ok := <-timer.C:
		case resp, ok := <-req.Chan:
			if !ok {
				return c.String(http.StatusInternalServerError, "channel closed or something")
			}
			return c.JSON(200, resp)
		case _ = <-timer.C:
			// timeout has elapsed
			resp := fmt.Sprintf("error: timeout reached (%s), at request to :%s", DefaultTimeout, c.Path())
			return c.String(400, resp)
		}
		return c.String(http.StatusInternalServerError, "FALSE")
	}
}

// server.GET("/api/numTabs", NumTabs)
func NumTabs(c echo.Context) error {
	start := time.Now()
	a := pwsh.RunScript("waterfoxTabs.ps1")
	fmt.Println("[Debug] (", c.Path(), ") request elapsed time:", time.Since(start))
	return c.JSON(200, a)
}

type TabData struct {
	Url      string `json:"url"`
	Title    string `json:"title"`
	Id       int    `json:"id"`
	WindowId int    `json:"windowId"`
	//Index int
}

// replace single asterisks with double asterisks
func replaceSingleAsterisks(s string) string {
	var result strings.Builder
	result.Grow(len(s)) // Pre-allocate for efficiency

	i := 0
	for i < len(s) {
		if s[i] == '*' {
			// Check if it's a single asterisk
			if i+1 < len(s) && s[i+1] == '*' {
				// It's a double asterisk, skip both
				result.WriteString("**")
				i += 2
			} else {
				// It's a single asterisk, replace with double
				result.WriteString("**")
				i++
			}
		} else {
			// Not an asterisk, just append the character
			result.WriteByte(s[i])
			i++
		}
	}

	return result.String()
}

// ~44ms at 145 tabs
func HandleFilteredGetTabs(w *wsu.WsConfig) echo.HandlerFunc {
	return func(c echo.Context) error {
		// TODO: remove
		c.Response().Header().Set("Access-Control-Allow-Origin", "*")

		// test if the webext is connected.
		if len(w.Connections) < 1 {
			resp := fmt.Sprintf("error=there are no connections to a browser extension. "+
				"first run the browser extension to enable retrieving browser data "+
				"at request to :%s", c.Path())
			return c.String(400, resp)
		}

		//simpleQuery := false
		var g glob.Glob
		filter := c.QueryParam("query")
		if filter == "" {
			return c.String(400, "no query passed.")
		}
		if strings.Contains(filter, "*") {
			filter = replaceSingleAsterisks(filter)
		} else {
			if strings.Count(filter, " ") == 0 {
				//simpleQuery = true
				filter = fmt.Sprintf("**%s**", filter)
			} else {
				parts := strings.Split(filter, " ")
				filter = strings.Join(parts, "**")
			}
		}
		//fmt.Println("filter=", filter)
		g, err := glob.Compile(filter, '.')
		if err != nil {
			return c.String(400, "HandleFilteredGetTabs: error while processing query filter:%s"+err.Error())
		}

		// create a timeout
		timer := time.NewTimer(3 * time.Second)
		//defer timer.Stop()

		//wsu.PendingWSRequestsV2.Mutex.Lock()
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
		w.Broadcast(bc)

		select {
		//case _, ok := <-timer.C:
		case resp, ok := <-req.Chan:
			//start := time.Now()
			if !ok {
				return c.String(http.StatusInternalServerError, "channel closed or something")
			}

			tabsData, ok := resp["response"]
			if !ok {
				return c.String(400, "HandleFilteredGetTabs: Unexpected error"+
					"while processing `getTabs` response from web-ext"+
					"err=`response` key not found in response object. this"+
					"may be some problem with the web extension.")
			}

			bytes, err := json.Marshal(tabsData)
			if err != nil {
				return c.String(400, err.Error())
			}

			var urls []TabData

			err = json.Unmarshal(bytes, &urls)
			if err != nil {
				return c.String(400, err.Error())
			}
			//fmt.Printf(`query data recieved from browser and prepared in: %s`, time.Since(start))

			/*
				simpleQuery := false
				if simpleQuery == true {
					var responseUrls []TabData
					for _, tab := range urls {
						if g.Match(strings.Join([]string{tab.Title, tab.Url}, "")) {
							responseUrls = append(responseUrls, tab)
							continue
						}

						// // 34.7 213 37 37 37 39 43 36
						// if strings.Contains(strings.Join([]string{tab.Title, tab.Url}, ""), filter) {
						// 	responseUrls = append(responseUrls, tab)
						// 	continue
						// }

						// // 45.87 46.1 39.3 39.3 43 39.0
						// if strings.Contains(tab.Title, filter) || strings.Contains(tab.Url, filter) {
						// 	responseUrls = append(responseUrls, tab)
						// 	continue
						// }
						// // 187.4 41.8 42.5 40.9 46.4
						// if strings.Contains(tab.Title, filter) {
						// 	responseUrls = append(responseUrls, tab)
						// 	continue
						// }
						// if strings.Contains(tab.Url, filter) {
						// 	//if g.Match(strings.ReplaceAll(tab.Url, ".", "*")) {
						// 	responseUrls = append(responseUrls, tab)
						// 	continue
						// }
					}
					return c.JSON(200, responseUrls)
				}
			*/

			var responseUrls []TabData
			for _, tab := range urls {
				if g.Match(tab.Title) {
					responseUrls = append(responseUrls, tab)
					continue
				}
				if g.Match(tab.Url) {
					//if g.Match(strings.ReplaceAll(tab.Url, ".", "*")) {
					responseUrls = append(responseUrls, tab)
					continue
				}
			}

			return c.JSON(200, responseUrls)
		case _ = <-timer.C:
			// timeout has elapsed
			resp := fmt.Sprintf("error: timeout reached (%s), at request to :%s", DefaultTimeout, c.Path())
			return c.String(400, resp)
		}
		return c.String(http.StatusInternalServerError, "FALSE")
	}
}

func HandleWebExtWS(c echo.Context, cfg *wsu.WsConfig) error {
	// this runs whenever a new connection request is made.

	// if c.QueryParam("context") == "webext" {...}

	// define a websocket handler.
	hndl := func(ws *websocket.Conn) {
		if len(cfg.Connections) >= 1 {
			fmt.Println("[WARN] You attempted to add another websocket connection, " +
				"but a previous connection(s) existed. " +
				"connections other than this new connection will be removed.")
			for ws := range cfg.Connections {
				ws.Close()
				cfg.RemoveConnection(ws)
			}
		}
		cfg.AddConnection(ws)
		defer func() {
			ws.Close()
			cfg.RemoveConnection(ws)
		}()

		handleV1Pended := func(id string, msg map[string]interface{}) {
			//PendingRequests.Remove(id) // make this only run once per id or something...
			//PendingRequests.Resolve(id, msg)
		}
		//#region read
		for {
			var msg []byte
			err := websocket.Message.Receive(ws, &msg)
			if err != nil {
				if errors.Is(err, io.EOF) {
					// fmt.Println("[LOG] <HandleWS> Client disconnected (EOF)")
					// cfg.RemoveConnection(ws)
					return // stop serving.
				} //else
				fmt.Println("[WARN] HandleWebExtWS: Unexpected client connection error:", err)
				// cfg.RemoveConnection(ws)
				return // stop serving.
			}

			// try to parse message as json. if the message is not json, discard the message.
			var jsonMsg map[string]interface{}
			err = json.Unmarshal(msg, &jsonMsg) //unmarshal into `jsonMsg`
			if err != nil {
				fmt.Println("ERROR:", err)
				// continue serving ...
			}
			id, ok := jsonMsg["id"].(string)
			if ok {
				if id != "" {
					PendingRequestsV2.Mutex.Lock()
					handleV1Pended(id, jsonMsg)
					for _, req := range PendingRequestsV2.Requests {
						if req.Id == id {
							//fmt.Println("[MATCH] recieved Request:", id)
							req.Chan <- jsonMsg // resolve the channel.
							break
						}
					}
					PendingRequestsV2.Mutex.Unlock()
				}
			}
			//fmt.Println("[LOG] <HandleWebExtWS> Received: ", jsonMsg)
			// func() {
			// 	a := string(msg)
			// 	if len(a) >= 255 {
			// 		fmt.Println("[LOG] <HandleWebExtWS> Received: ", string(a[0:254]))
			// 	} else {
			// 		fmt.Println("[LOG] <HandleWebExtWS> Received: ", string(a))
			// 	}

			// }()

			// if the message has an id field, process it.

			//fmt.Println("client says", msg)

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
