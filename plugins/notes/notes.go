package notes

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/eissar/nest/config"
	"github.com/eissar/nest/plugins/pwsh"
	"github.com/labstack/echo"
)

var Editor = "C:/Program Files/Neovim/bin/nvim.exe"

// server.GET("/api/recentNotes", core.RecentNotes)

func RecentNotes(c echo.Context) error {
	a := pwsh.RunScript("recentNotes.ps1")
	return c.JSON(200, a)
}

func Edit(c echo.Context) error {
	start := time.Now()
	leaf := c.QueryParam("fileName") // Get-ChildItem "$env:CLOUD_DIR/Catalog/*.md"
	fp := path.Join(os.Getenv("CLOUD_DIR"), "Catalog", leaf)
	fmt.Println(Editor, fp)
	//fp = "C:/Users/eshaa/draft.lua"
	open_cmd := fmt.Sprintf("wt.exe -d $env:CLOUD_DIR nvim %v", fp)
	pwsh.Exec(open_cmd)
	fmt.Println("[Debug] (", c.Path(), ") request elapsed time:", time.Since(start))
	return c.String(http.StatusOK, "success")
}

/*
server.GET("/api/recentNotes", core.RecentNotes)
server.POST("/api/edit", core.Edit)
*/
//func PopulateGetNotesDetail(c echo.Context, templateName string) interface{} {
//	return GetNotesNamesDates(-1, 0)
//}

func RegisterRootRoutes(server *echo.Echo) {
	server.POST("/api/edit", Edit)
	server.GET("/api/recentNotes", RecentNotes)

}
