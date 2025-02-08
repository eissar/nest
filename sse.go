package main

import (
	"fmt"
	"net/http"
	"sync"
	. "web-dashboard/types"

	"github.com/labstack/echo/v4"
)

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
	fmt.Println("[LOG] <HandleSSE> Client disconnected")
	return nil
}

func SendSSEAll(message string) {
	clientsMu.RLock()
	defer clientsMu.RUnlock()
	for client := range clients {
		select { // Non-blocking send to avoid deadlocks
		case client.Ch <- message:
		default:
			// If the client's channel is full, it means the client
			// might have disconnected.  We don't want to block here.
			fmt.Println("Client's channel full. Skipping message.")
		}

	}
}
func SendSSETargeted(message string, clients map[*Client]bool) {
	clientsMu.RLock()
	defer clientsMu.RUnlock()
	for client := range clients {
		select { // Non-blocking send to avoid deadlocks
		case client.Ch <- message:
		default:
			// If the client's channel is full, it means the client
			// might have disconnected.  We don't want to block here.
			fmt.Println("Client's channel full. Skipping message.")
		}

	}
}

// api function
// inputs:
// c echo.Context
// kind string : type of broadcast for case matching.
func BroadcastEvent(c echo.Context, kind string) error {
	fmt.Println("[LOG] <BroadcastEvent> sending event", kind)
	sent := false
	switch kind {
	case "getSong":
		SendSSEAll(`{"event":"getSong"}`)
		sent = true
	case "ytMusicElement":
		elem := c.QueryParam("elem")
		var message string
		// TODO: Create case to take songName param

		if len(elem) >= 1 {
			message = elem
		}
		if !(len(message) >= 1) {
			return c.String(400, "no query param.")
		}
		Broadcast(message)
		sent = true
	}
	if sent {
		return c.String(http.StatusOK, "OK")
	} else {
		return c.String(400, "server err")
	}
}

// func RetrieveDataAndExecuteTemplate(c echo.Context, name string, data any) {
// 	a := GetEnumerateWindows()[0]
// 	err = templ.ExecuteTemplate(c.Response().Writer, "window.html", a)
// 	return nil
// }
