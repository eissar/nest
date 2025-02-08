package apiroutes

import (
	"fmt"
	"sync"
	. "web-dashboard/types"
)

func SendSSETargeted(message string, clients map[*Client]bool, mu *sync.RWMutex) {
	mu.RLock()
	defer mu.RUnlock()
	fmt.Println("[DEBUG] sending message", message)
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
