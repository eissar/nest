package browser

// go get -u github.com/eissar/browser-query@master
import (
	"fmt"
	"sync"

	"github.com/eissar/browser-query"
	"github.com/eissar/nest/config"
	"github.com/labstack/echo/v4"
)

// TODO: ? merge from github.com/eissar/browser-query
// WARN: browserquery has state, bad library; reformat before using here.

// var (
// 	clientsMu sync.RWMutex
// 	clients   = make(map[*browserQuery.Client]bool)
// )

func RegisterRootRoutes(n config.NestConfig, server *echo.Echo) {
	// server.GET("/eagleApp/sse", browserQuery.HandleSSE)
	// server.GET("/api/uploadTabs", browserQuery.HandleSSE)

	// TEST
	// if nestConfig.Nest.Plugins.browser == true {
	//
	server.GET("/eagleApp/sse", browserQuery.HandleSSE)
	//server.POST("/api/uploadTabs", browserQuery.UploadTabs)

	server.POST(
		"/api/uploadTabs",
		browserQuery.UploadTabsHandler(func(c echo.Context, t []browserQuery.TabInfo) {
			fmt.Println("tabs count:", len(t)) // works
		}),
	)
	// }

	// TEST

}
