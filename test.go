package main

import (
	"fmt"

	//handlers "github.com/eissar/nest/handlers"

	"github.com/labstack/echo/v4"
)

func RegisterTestRoutes(g *echo.Group) {
	//g.GET("/test", handlers.DynamicTemplateHandler("notes-struct.html", notes.PopulateGetNotesDetail))

	// test channels
	g.GET("api/chan", func(c echo.Context) error {
		chanS := struct {
			Ch chan string
		}{
			Ch: make(chan string),
		}
		go func() {
			chanS.Ch <- "1"
		}()
		v := <-chanS.Ch
		fmt.Printf("v: %v\n", v)

		return c.String(200, "ok")
	})
}
