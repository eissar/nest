package main

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"golang.org/x/net/websocket"
)

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
			broadcast(broadcastMessage)
		}
	}).ServeHTTP(c.Response(), c.Request())
	return nil
}
