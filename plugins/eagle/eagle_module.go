package eaglemodule

import (
	"fmt"
	"sync"

	"github.com/eissar/nest/config"
	"github.com/eissar/nest/helpers"
	"github.com/eissar/nest/plugins/pwsh"

	"github.com/labstack/echo/v4"
)

type Library struct {
	Name  string
	Path  string
	Mutex sync.Mutex
}
type Eagle struct {
	//Db        *sql.DB
	Libraries []*Library
}

func RegisterGroupRoutes(g *echo.Group) {
	//nestCfg := GetConfig()

	g.GET("/testcfg", func(c echo.Context) error {
		config.MustNewConfig()
		return c.JSON(200, "OK")
	})
	g.GET("/getcfg", func(c echo.Context) error {
		return c.JSON(
			200,
			config.GetConfig(),
		)
	})
	// TODO: MAKE WORK
	g.GET("/isValid/:id", func(c echo.Context) error {
		id := c.Param("id")
		if id == "" {
			return c.JSON(200, `{"valid":false}`)
		}
		return c.JSON(200, `{"valid":true}`)
	})
	g.GET("/test", func(c echo.Context) error {
		a, err := validateIsEagleServerRunning("http://localhost:41595/api/application/info")
		if err != nil {
			return c.String(400, "NO")
		}
		return c.JSON(200, a)
	})

	g.GET("/open/:id", func(c echo.Context) error {
		id := c.Param("id")
		uri := fmt.Sprintf("eagle://item/%s", id)
		helpers.OpenURI(uri)
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
	server.GET("/eagle\\://item/:itemId", ServeThumbnailHandler(&n))
	server.GET("/:itemId", ServeThumbnailHandler(&n))
	server.GET("/api/eagleOpen/:id", func(c echo.Context) error {
		id := c.Param("id")
		uri := fmt.Sprintf("eagle://item/%s", id)
		helpers.OpenURI(uri)
		return c.String(200, "OK")
	})

	/*
		server.GET("/:id", func(c echo.Context) error {
			id := EagleItemId(c.Param("id"))
			if !id.IsValid() {
				res := fmt.Sprintf("get path=%s err=id of `%s` is not valid.", c.Path(), id)
				return c.JSON(200, res)
			}

			thumbnail, err := nestCfg.GetEagleThumbnail(id)
			if err != nil {
				res := fmt.Sprintf("get path=%s err=%s", c.Path(), err.Error())
				return c.String(200, res)
			}
			err = fileUtils.PathExists(thumbnail.ThumbnailPath)
			if err != nil {
				res := fmt.Sprintf("get path=%s err=%s", c.Path(), err.Error())
				return c.String(200, res)
			}
			// filepath exists.
			return c.File(thumbnail.ThumbnailPath)
			//return c.JSON(200, thumb)
		})
	*/
}
