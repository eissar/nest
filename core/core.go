package core

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

func Ping(c echo.Context) error {
	return c.String(http.StatusOK, "OK")
}

// func PopulateQueryFrontmatter() {}
// func PopulateGetRecentNotes()   {}

func ServerShutdown(c echo.Context) error {
	var err error

	//use go routine with timeout to allow time for response.
	timeout := 10 * time.Second
	timeoutCtx, shutdownRelease := context.WithTimeout(context.Background(), timeout)
	defer shutdownRelease()

	go func() {
		err = c.Echo().Server.Shutdown(timeoutCtx)
	}()
	if err != nil {
		fmt.Println("err while graceful shutdown:", err)
	}

	return c.String(200, "shutdown cmd successful.")
}
func Shutdown(s *echo.Echo) error {
	var err error
	//use go routine with timeout to allow time for response.
	timeout := 10 * time.Second
	timeoutCtx, shutdownRelease := context.WithTimeout(context.Background(), timeout)
	defer shutdownRelease()

	go func() {
		err = s.Shutdown(timeoutCtx)
	}()
	if err != nil {
		fmt.Println("err while graceful shutdown:", err)
	}
	return nil
}

type uploadTabsBody struct {
	Body string `json:"body"`
}

func UploadTabs(c echo.Context) error {
	a := c.Request().Body
	b := []byte{}
	_, err := a.Read(b)
	if err != nil {
		panic(err)
	}
	fmt.Println("[SUCCESS]", c)
	return c.String(200, "OK")
}

// registers routes on the server root (/)
func RegisterRootRoutes(server *echo.Echo) {
	//server.GET("/eagle\\://item/:itemId", ServeThumbnailHandler(&n))
	server.GET("/api/server/close", ServerShutdown)
	server.GET("/api/ping", Ping)
	server.POST("/api/uploadTabs", UploadTabs)

}
