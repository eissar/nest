package ytmmodule

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"io"
	"sync"
	"time"
	wsu "web-dashboard/websocket-utils"

	"github.com/labstack/echo/v4"
	"golang.org/x/net/websocket"
)

var lastSong = "NULL SONG DATA"

func HandleWS(c echo.Context, cfg *wsu.WsConfig) error {
	websocket.Handler(func(conn *websocket.Conn) {

		cfg.AddConnection(conn)
		defer func() {
			cfg.RemoveConnection(conn)
			conn.Close()
		}()

		// this is used to both serve data to the dom
		// and receive data updates from browser extension/
		// userscripts.

		// which context we are can be retrieved from the
		// c.QueryParam("context")

		// Initial message (dependent on context)

		pingTicker := time.NewTicker(30 * time.Second)
		defer pingTicker.Stop()

		failsCnt := 0
		go func() {
			for range pingTicker.C {
				// Send a Ping message
				if err := websocket.Message.Send(conn, "ping"); err != nil {
					fmt.Println("[LOG] <HandleWS> Failed to send ping:", err)
					failsCnt += 1
					if failsCnt == 3 {
						// Exit goroutine if ping fails
						return
					}
				}
			}
		}()

		(func() {
			var initial_msg string
			if c.QueryParam("context") == "" {
				return
			}
			if c.QueryParam("context") == "webpage" {
				initial_msg = fmt.Sprintf(`<a id="message-container">%s</a>`, lastSong)
			}
			if c.QueryParam("context") == "webext" {
				fmt.Println("init message webext")
				initial_msg = fmt.Sprintf(`{"event":"getSong"}`, lastSong)
			}
			if err := websocket.Message.Send(conn, initial_msg); err != nil {
				fmt.Println("[ERROR] <initialMessage>", err)
				return
			}
		})()

		for {
			// Read
			var msg map[string]interface{}
			err := websocket.JSON.Receive(conn, &msg)
			if err != nil {
				if errors.Is(err, io.EOF) {
					fmt.Println("[LOG] <HandleWS> Client disconnected (EOF)")
					return
				} else {
					fmt.Println("[ERROR] <HandleWS> Unexpected error receiving message:", err)
					return
				}
			}

			fmt.Println("Received:", msg)

			// why this?
			if cmd, ok := msg["command"].(string); ok {
				if cmd == "getSong" {
					cfg.Broadcast(`{"event":"getSong"}`)
				}
			}

			// broadcastMessage := fmt.Sprintf("Client says: %s", msg)
			// cfg.Broadcast(broadcastMessage)
		}
	}).ServeHTTP(c.Response(), c.Request())
	return nil
}

var WSConfig wsu.WsConfig

type songData struct {
	Song     string `json:"song"`
	Duration string `json:"duration"`
	Channel  string `json:"channelName"`
}

/*
   params.append("song", playingSong);
   params.append("album", getAlbumName());
   params.append("channelName", getChannelName());
   params.append("channelUrl", getChannelUrl());
   params.append("duration", getSongDuration());
   params.append("watchId", getWatchId());
*/

// ytm.RegisterRoutesFromGroup
// prefix ytm
func RegisterRoutesFromGroup(g *echo.Group) {

	templ := template.Must(template.ParseFiles("templates/song.templ"))

	WSConfig := wsu.WsConfig{
		Connections: make(map[*websocket.Conn]bool),
		Mutex:       sync.Mutex{},
	}
	g.GET("/ws", func(c echo.Context) error {
		err := HandleWS(c, &WSConfig)
		if err != nil {
			fmt.Println("[ERROR]", err.Error())
			c.String(400, err.Error())
		}
		return c.String(200, "OK")
	})

	//#region broadcasts to websockets.
	g.GET("/broadcast/getSong", func(c echo.Context) error {
		WSConfig.Broadcast(`{"event":"getSong"}`)
		return c.String(200, "OK")
	})
	// e.g, /playingSong?song=<songName>
	g.POST("/playingSong", func(c echo.Context) error {
		lastSong = c.QueryParam("song")
		if lastSong == "" {
			return c.String(400, "OK")
		}

		dat := songData{
			Song:     c.QueryParam("song"),
			Channel:  c.QueryParam("channelName"),
			Duration: c.QueryParam("duration"),
		}
		//a := &string
		var buf bytes.Buffer
		templ.ExecuteTemplate(&buf, "song.templ", dat)
		// send an element inline will autoupdate dom with
		// matching id.

		out := buf.String()
		fmt.Println("data:", out)
		WSConfig.Broadcast(out)
		return c.String(200, "OK")
	})
	g.GET("/playingSong", func(c echo.Context) error {
		lastSong = c.QueryParam("song")
		// send an element inline will autoupdate dom with
		// matching id.
		b := fmt.Sprintf(`<a id="message-container">%s</a>`, lastSong)
		WSConfig.Broadcast(b)
		return c.String(200, "OK")
	})
	//#endregion

	// on route registration, send a broadcast.
	WSConfig.Broadcast(`{"event":"getSong"}`)
}
