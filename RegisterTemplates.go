package main

import (
	"github.com/labstack/echo/v4"
	//	apiroutes "github.com/eissar/nest/api-routes"
	handlers "github.com/eissar/nest/handlers"
)

func RegisterTemplateRoutes(g *echo.Group) {
	//g.GET("/notes-struct", handlers.DynamicTemplateHandler("notes-struct.html", apiroutes.PopulateGetNotesDetail))
	//g.GET("/windows", handlers.DynamicTemplateHandler("windows.html", apiroutes.PopulateEnumerateWindows))
	g.GET("/recent-notes", handlers.PwshTemplateHandler("recent-notes.html", "recentNotes.ps1"))
	g.GET("/key-value", handlers.PwshTemplateHandler("key-value.templ", "mock_nvim.ps1"))
	g.GET("/open-tabs-count", handlers.PwshTemplateHandler("open-tabs-count.templ", "waterfoxTabs.ps1"))
	//server.GET("/template/open-tabs", dynamicTemplateHandler("open-tabs.templ", apiroutes.PopulateOpenTabs))
	g.GET("/recent-eagle-items", handlers.PwshTemplateHandler("recent-eagle-items.templ", "recentEagleItems.ps1"))

	g.GET("/sse-browser-tabs", handlers.StaticTemplateHandler("sse-browser-tabs.templ"))
	g.GET("/browser-tabs", handlers.StaticTemplateHandler("browser-tabs.templ"))
	g.GET("/recent-notes_layout", handlers.StaticTemplateHandler("recent-notes.layout.html"))
	g.GET("/timeline_layout", handlers.StaticTemplateHandler("timeline.layout.html"))
	g.GET("/now-playing", handlers.StaticTemplateHandler("ws-now-playing.ytm.templ")) // ./templates/ws-now-playing.ytm.templ
	g.GET("/ping", handlers.StaticTemplateHandler("ping.templ"))
}
