package main

import (
	"fmt"
	"sync"

	"github.com/labstack/echo/v4"
	"golang.org/x/net/websocket"
)

// Connection Management
var (
	connections = make(map[*websocket.Conn]bool)
	mutex       sync.Mutex
)

func addConnection(ws *websocket.Conn) {
	mutex.Lock()
	defer mutex.Unlock()
	connections[ws] = true
}

func removeConnection(ws *websocket.Conn) {
	mutex.Lock()
	defer mutex.Unlock()
	delete(connections, ws)
}

func Broadcast(message string) {
	mutex.Lock()
	defer mutex.Unlock()
	lastSong = message
	for ws := range connections {
		go func(ws *websocket.Conn) {
			// Send concurrently to avoid blocking
			if err := websocket.Message.Send(ws, message); err != nil {
				fmt.Println("Error sending to a connection:", err)
				removeConnection(ws) // Remove the connection if sending fails
				ws.Close()
			}
		}(ws)
	}
}
func Hello(c echo.Context) error {
	websocket.Handler(func(ws *websocket.Conn) {
		addConnection(ws)
		defer func() {
			removeConnection(ws)
			ws.Close()
		}()

		// Initial message (optional)
		initial_msg := fmt.Sprintf(`<a id="message-container">%s</a>`, lastSong)
		if err := websocket.Message.Send(ws, initial_msg); err != nil {
			fmt.Println(err)
			return
		}

		for {
			// Read
			msg := ""
			err := websocket.Message.Receive(ws, &msg)
			if err != nil {
				// ... (Handle errors as before: EOF, network issues, etc.) ...
				if err.Error() == "EOF" {
					fmt.Println("Client disconnected (EOF)")
					return
				} else {
					fmt.Println("Unexpected error receiving message:", err)
				}
			}

			fmt.Printf("Received: %s\n", msg)

			// Broadcast the received message to all clients
			broadcastMessage := fmt.Sprintf("Client says: %s", msg)
			Broadcast(broadcastMessage)
		}
	}).ServeHTTP(c.Response(), c.Request())
	return nil
}
