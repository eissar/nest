package nest

import (
	"encoding/json"
	"fmt"
	"log"
	//"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/eissar/nest/api"
	"github.com/eissar/nest/config"
	"github.com/eissar/nest/plugins/launch"
	"github.com/eissar/nest/plugins/pwsh"
	"github.com/fsnotify/fsnotify"
	"golang.org/x/net/context"

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

// wait until library is switched
func pollForSwitch(c context.Context, pollingCh chan bool, targetLib string) {
	t := time.NewTicker(500 * time.Millisecond)

	for {
		select {
		case <-t.C: // tick
			currLib, err := CurrentLibraryPath()
			if err != nil {
				log.Printf("warning err=%v\n", err)
				continue
			}
			if currLib == targetLib {
				// library has switched
				pollingCh <- true
			} else {
				fmt.Printf("currLib: %v\n", currLib)
				fmt.Printf("targetLib: %v\n", targetLib)
			}

		case <-c.Done(): // exit
			return
		}
	}
}

func LibrarySwitchSync(baseUrl string, libraryPath string) error {
	currLibraryPath, err := CurrentLibraryPath()
	if err != nil {
		return fmt.Errorf("libraryswitchsync: error getting lib info err=%v", err)
	}

	fmt.Printf("libraryPath: %v\n", libraryPath)
	fmt.Printf("currLibraryPath: %v\n", currLibraryPath)

	if currLibraryPath == libraryPath {
		// do nothing.
		return nil
	}

	// switch library
	timeoutCh := time.After(10 * time.Second)

	err = api.SwitchLibrary(baseUrl, libraryPath)
	if err != nil {
		return fmt.Errorf("couldn't switch to lib=%s err=%v", libraryPath, err)
	}

	ctx, cancelPolling := context.WithCancel(context.Background())
	pollingCh := make(chan bool)
	go pollForSwitch(ctx, pollingCh, libraryPath)
	defer cancelPolling()

	for {
		select {
		case <-pollingCh:
			fmt.Println("SWITCHED")
			return nil
		case <-timeoutCh:
			return fmt.Errorf("timeout elapsed")
		}
	}
}

// lol idiot
func CurrentLibrary0() (string, error) {
	cfg := config.GetConfig()

	// 1. get first item id

	resp, err := api.ListV3(cfg.BaseURL(), 1)
	if err != nil {
		return "", fmt.Errorf("currentlibrary: could not retrieve any library items.")
	}

	firstItemId := resp.Data[0].Id
	if firstItemId == "" {
		return "", fmt.Errorf("currentlibrary: could not retrieve any library items.")
	}

	// 2. get thumbnail

	thumb, err := api.Thumbnail(cfg.BaseURL(), firstItemId)
	if err != nil {
		return "", fmt.Errorf("could not find thumbnail for first eagle item err=%v", err)
	}

	//fmt.Println("thumb:", thumb)
	parts := strings.SplitAfterN(thumb, `.library/`, 2)
	if len(parts) == 1 {
		return "", fmt.Errorf("currentlibrary: could not segment the first retrieved eagle item path by .library does the eagle library have some other name?")
	}
	currLib := filepath.Clean(parts[0])

	// optionally check for metadata.json?
	return currLib, nil
}

func CurrentLibrary() (*api.Library, error) {
	cfg := config.GetConfig()
	libInfo, err := api.GetLibraryInfo(cfg.BaseURL())
	if err != nil {
		return nil, fmt.Errorf("error getting library info err=%w", err)
	}
	//fmt.Println(libInfo.Data.Library)
	return &libInfo.Data.Library, nil
}

// get path to current library
func CurrentLibraryPath() (string, error) {
	currLib, err := CurrentLibrary()
	if err != nil {
		return "", fmt.Errorf("error getting library info err=%w", err)
	}
	return currLib.Path, nil
}

func CurrentLibraryName() (string, error) {
	cfg := config.GetConfig()
	libInfo, err := api.GetLibraryInfo(cfg.BaseURL())
	if err != nil {
		return "", fmt.Errorf("error getting library info err=%w", err)
	}

	return libInfo.Status, nil
	//return libInfo.Data.Library
}
