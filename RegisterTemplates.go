package main

import (
	"github.com/labstack/echo/v4"
	//	apiroutes "web-dashboard/api-routes"
	handlers "web-dashboard/handlers"
	pwsh "web-dashboard/powershell-utils"
)

func RegisterTemplateRoutes(g *echo.Group) {
	//g.GET("/notes-struct", handlers.DynamicTemplateHandler("notes-struct.html", apiroutes.PopulateGetNotesDetail))
	//g.GET("/windows", handlers.DynamicTemplateHandler("windows.html", apiroutes.PopulateEnumerateWindows))
	g.GET("/recent-notes", handlers.PwshTemplateHandler("recent-notes.html", pwsh.PwshScript, "./powershell-utils/recentNotes.ps1"))
	g.GET("/key-value", handlers.PwshTemplateHandler("key-value.templ", pwsh.PwshScript, "./powershell-utils/mock_nvim.ps1"))
	g.GET("/open-tabs-count", handlers.PwshTemplateHandler("open-tabs-count.templ", pwsh.PwshScript, "./powershell-utils/waterfoxTabs.ps1"))
	//server.GET("/template/open-tabs", dynamicTemplateHandler("open-tabs.templ", apiroutes.PopulateOpenTabs))
	g.GET("/recent-eagle-items", handlers.PwshTemplateHandler("recent-eagle-items.templ", pwsh.PwshScript, "./powershell-utils/recentEagleItems.ps1"))
	g.GET("/sse-browser-tabs", handlers.StaticTemplateHandler("sse-browser-tabs.templ"))
	g.GET("/browser-tabs", handlers.StaticTemplateHandler("browser-tabs.templ"))
	g.GET("/recent-notes_layout", handlers.StaticTemplateHandler("recent-notes.layout.html"))
	g.GET("/timeline_layout", handlers.StaticTemplateHandler("timeline.layout.html"))
	g.GET("/now-playing", handlers.StaticTemplateHandler("ws-now-playing.ytm.templ")) // ./templates/ws-now-playing.ytm.templ
	g.GET("/ping", handlers.StaticTemplateHandler("ping.templ"))
}
