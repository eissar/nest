package websocketutils

import (
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"golang.org/x/net/websocket"
)

type WsConfig struct {
	Connections map[*websocket.Conn]bool `json:"connections"`
	Mutex       sync.Mutex
}

// c *websocket.Conn connection to add
func (cfg *WsConfig) AddConnection(c *websocket.Conn) {
	cfg.Mutex.Lock()
	defer cfg.Mutex.Unlock()
	cfg.Connections[c] = true
}

// c *websocket.Conn connection to remove
func (cfg *WsConfig) RemoveConnection(c *websocket.Conn) {
	cfg.Mutex.Lock()
	defer cfg.Mutex.Unlock()
	delete(cfg.Connections, c)
}

func WSBroadcastTargeted(c echo.Context, ws *WsConfig) error {
	return nil
}

type PendingWSRequest struct {
	Start time.Time `json:"start"`
	Id    string    `json:"id"`
}

type PendingWSRequests struct {
	Requests []*PendingWSRequest
	Mutex    sync.Mutex
}

// var
// var pendingRequests pendingWSRequestsCfg{
// 	Requests: []pendingWSRequest{},
// 	Mutex: sync.Mutex{}
// }

func (wsReq *PendingWSRequests) New() *PendingWSRequest {
	wsReq.Mutex.Lock()
	defer wsReq.Mutex.Unlock()

	// create a pointer to a new pending request
	a := &PendingWSRequest{
		Start: time.Now(),
		Id:    uuid.NewString(),
	}
	wsReq.Requests = append(wsReq.Requests, a)
	return a
}

// remove any Pending requests older than the timeout
func (wsReq *PendingWSRequests) Trim() {
	wsReq.Mutex.Lock()
	defer wsReq.Mutex.Unlock()

	t := 1 * time.Minute

	newArr := []*PendingWSRequest{}
	for _, req := range wsReq.Requests {
		if time.Since(req.Start) >= t {
			return
		}
		// otherwise, add it to newArr
		newArr = append(newArr, req)
	}
	wsReq.Requests = newArr
}

func (ws *WsConfig) Broadcast(message string) {
	ws.Mutex.Lock()
	defer ws.Mutex.Unlock()

	for conns := range ws.Connections {
		go func(conn *websocket.Conn) {
			// Send concurrently to avoid blocking
			if err := websocket.Message.Send(conn, message); err != nil {
				fmt.Println("Error sending to a connection:", err)
				ws.RemoveConnection(conn) // Remove the connection if sending fails
				conn.Close()
			}
		}(conns)
	}
}
