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
	Chan        *chan int
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

type WSRequest struct {
	Start time.Time `json:"start"`
	Id    string    `json:"id"`
	Chan  chan map[string]interface{}
}
type WSRequestV2 struct {
	Start time.Time `json:"start"`
	Id    string    `json:"id"`
	Chan  chan map[string]interface{}
}

type PendingWSRequests struct {
	Requests []*WSRequest
	Mutex    sync.Mutex
}
type PendingWSRequestsV2 struct {
	Requests []*WSRequestV2
	Mutex    sync.Mutex
}

// var
// var pendingRequests pendingWSRequestsCfg{
// 	Requests: []pendingWSRequest{},
// 	Mutex: sync.Mutex{}
// }

// adds a request to wsReq and returns a pointer to the WSRequest
func (wsReq *PendingWSRequests) New() *WSRequest {
	wsReq.Mutex.Lock()
	defer wsReq.Mutex.Unlock()

	// create a pointer to a new pending request
	a := &WSRequest{
		Start: time.Now(),
		Id:    uuid.NewString(),
	}
	//fmt.Println("[DEBUG] <PendingWsRequests.New> websocket request with ID: ", a.Id, "created.")
	wsReq.Requests = append(wsReq.Requests, a)
	return a
}

// adds a request to wsReq and returns a pointer to the WSRequest
func (wsReq *PendingWSRequestsV2) New() *WSRequestV2 {
	wsReq.Mutex.Lock()
	defer wsReq.Mutex.Unlock()

	// create a pointer to a new pending request
	a := &WSRequestV2{
		Start: time.Now(),
		Id:    uuid.NewString(),
		Chan:  make(chan map[string]interface{}),
	}
	//fmt.Println("[DEBUG] <PendingWsRequests.New> websocket request with ID: ", a.Id, "created.")
	wsReq.Requests = append(wsReq.Requests, a)
	return a
}

// remove any Pending requests older than the timeout
// timeout: default 1min
func (wsReq *PendingWSRequests) Trim() {
	wsReq.Mutex.Lock()
	defer wsReq.Mutex.Unlock()

	t := 1 * time.Minute

	newArr := []*WSRequest{}
	for _, req := range wsReq.Requests {
		if time.Since(req.Start) >= t {
			return
		}
		// otherwise, add it to newArr
		newArr = append(newArr, req)
	}
	wsReq.Requests = newArr
}

func (wsReq *PendingWSRequests) Remove(id string) {
	wsReq.Mutex.Lock()
	defer wsReq.Mutex.Unlock()

	var array_requests []*WSRequest
	for _, r := range wsReq.Requests {
		if r.Id == id {
			//fmt.Println("[DEBUG] <PendingWsRequests.Remove> websocket request ID: ", id, "resolved in: ", time.Since(r.Start))
			continue
		}
		array_requests = append(array_requests, r)
	}
	wsReq.Requests = array_requests
}
func (wsReq *PendingWSRequests) RemovePointer(ptrToRemove *WSRequest) {
	wsReq.Mutex.Lock()
	defer wsReq.Mutex.Unlock()

	for i, ptr := range wsReq.Requests {
		if ptr == ptrToRemove {
			l := len(wsReq.Requests)

			// replace ptr with the last element
			wsReq.Requests[i] = wsReq.Requests[(l - 1)]

			// trim off the last element
			wsReq.Requests = wsReq.Requests[:(l - 1)]
		}
	}
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

func (wsReq *PendingWSRequests) Resolve(id string, data any) {
	wsReq.Mutex.Lock()
	defer wsReq.Mutex.Unlock()
}
