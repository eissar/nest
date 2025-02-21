package main

import (
	"fmt"

	//handlers "github.com/eissar/nest/handlers"

	_ "github.com/eissar/nest/eagle/api"
	"github.com/labstack/echo/v4"
)

func test(c echo.Context) error {
	return c.String(200, "YAYA")
	//a, err := api.InvokeRaindropAPI(api.Endpoint{}, nil)
	//if err != nil {
	//	return c.String(400, fmt.Sprintf("err=%v", err))
	//}
	//return c.String(200, a)
}

func RegisterTestRoutes(g *echo.Group) {
	//g.GET("/test", handlers.DynamicTemplateHandler("notes-struct.html", notes.PopulateGetNotesDetail))

	g.GET("/a", test)
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
