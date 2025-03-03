package nest

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/eissar/nest/api"
	"github.com/eissar/nest/config"
	"github.com/eissar/nest/plugins/launch"
	"github.com/eissar/nest/plugins/pwsh"
	"github.com/fsnotify/fsnotify"

	"github.com/labstack/echo/v4"
)

// maybe
type Data interface {
	GetData() []any
}

type Library struct {
	Name  string
	Path  string
	Mutex sync.Mutex
}

// rename to nest?
type Eagle struct {
	//Db        *sql.DB
	Libraries []*Library
}

// for endpoints that return a string.
type EagleDataMessage struct {
	Status string `json:"status"`
	Data   string `json:"data"`
}

type validResponse struct {
	Valid bool `json:"valid,omitempty"`
}

func RegisterGroupRoutes(g *echo.Group) {
	//nestCfg := GetConfig()

	// @Summary      refresh config
	// @Description  refresh config
	// @Router       /{group}/testcfg [get]
	g.GET("/testcfg", func(c echo.Context) error {
		config.MustNewConfig()
		return c.JSON(200, "OK")
	})
	// @Summary      get config
	// @Router       /getcfg [get]
	// @Produce		application/json
	// @Success 200	{object} config.NestConfig
	g.GET("/getcfg", func(c echo.Context) error {
		return c.JSON(
			200,
			config.GetConfig(),
		)
	})
	// @Summary     is valid
	// @Param				id	path	string	true	"id to check"
	// @Produce			application/json
	// @Success			200	{object} validResponse
	// @Router      /isValid/{id} [get]
	g.GET("/isValid/:id", func(c echo.Context) error {
		id := c.Param("id")

		if api.IsValidItemID(id) {
			return c.JSON(200, validResponse{Valid: true})
		}
		return c.JSON(200, validResponse{Valid: false})
	})
	// @Summary     is eagle server running
	// @Router      /test [get]
	// @Success			200 {string} string
	g.GET("/test", func(c echo.Context) error {
		a, err := validateIsEagleServerRunning("http://localhost:41595/api/application/info")
		if err != nil {
			return c.String(200, "false")
		}
		return c.String(200, fmt.Sprintf("%v", a))
	})

	// @Summary     reveal Id in eagle.
	// @Router      /open/:id [get]
	// @Param				id	path	string	true	"id to reveal"
	// @Success			200 {string} OK
	g.GET("/open/:id", func(c echo.Context) error {
		id := c.Param("id")
		uri := fmt.Sprintf("eagle://item/%s", id)
		launch.OpenURI(uri)
		return c.String(200, "OK")
	})
	g.GET("/api/recentEagleItems", pwsh.RecentEagleItems)
	g.GET("/api/recent-items", pwsh.RecentEagleItems)
	g.GET("/recent-items", pwsh.RecentEagleItems)

	//g.GET("/eagle\\://item/:itemId", func(c echo.Context) error {
	//	fmt.Print("test")
	//	return c.String(200, "works")
	//})
	// g.GET("/eagle://item/:itemId", func(c echo.Context) error {

	// 	return c.String(200, "works")
	// })
}

// registers routes on the server root (/)
func RegisterRootRoutes(n config.NestConfig, server *echo.Echo) {
	// @Summary     serve image
	// @Router      /eagle://item/{id} [get]
	// @Router      /{id} [get]
	// @Param				id	path	string	true	"id to serve image"
	// @Param				fq	query	string	false	"flag for full-quality response"
	// @Produce  		image/png
	// @Success			200 {file} thumbnail
	server.GET("/eagle\\://item/:itemId", ServeThumbnailHandler(&n))
	server.GET("/:itemId", ServeThumbnailHandler(&n))

	// show the item in eagle
	server.GET("/eagle/show/:id", func(c echo.Context) error {
		id := c.Param("id")
		uri := fmt.Sprintf("eagle://item/%s", id)
		launch.OpenURI(uri)
		return c.String(200, "OK")
	})

	// todo make better
	server.GET("/eagle/item/path", func(c echo.Context) error {
		id := EagleItemId(c.QueryParam("id"))
		if !id.IsValid() {
			return c.String(404, "invalid or missing query param `id`")
		}

		return c.String(200, "OK")
	})

	//$ext = (irm "http://localhost:41595/api/item/info?id=M6VK3GF6845SQ").data.ext
}

//#region monitor mtime

type Mtime map[string]int

func TryIngestMtime(n config.NestConfig) *Mtime {
	mtime := filepath.Join(n.Libraries.Paths[0], "mtime.json")
	if _, err := os.Stat(mtime); err != nil {
		log.Fatalf("%s", err.Error())
	}

	out := &Mtime{}

	bytes, err := os.ReadFile(mtime)
	if err != nil {
		log.Fatalf("%s", err.Error())
	}

	err = json.Unmarshal(bytes, out)
	if err != nil {
		log.Fatalf("%s", err.Error())
	}

	return out
}

var state = &Mtime{}

func WatchMtime(n config.NestConfig) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	mtime := filepath.Join(n.Libraries.Paths[0], "mtime.json")
	if _, err = os.Stat(mtime); err != nil {
		log.Fatalf("%s", err.Error())
	}

	watcher.Add(mtime)

	// WARN: from my understanding mtime.json can get quite large.

	// TODO: optimize pointer use.
	done := make(chan bool)
	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			fmt.Println("event:", event)
			newState := *TryIngestMtime(n)
			lastState := *state
			for id, _time := range newState {
				if lastState[id] != newState[id] {
					// this is the item that has changed.
					fmt.Println(id, _time)
				}
			}
			state = &newState

		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			log.Println("error:", err)
		}
	}
	<-done
}

//#endregion
