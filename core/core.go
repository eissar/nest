package core

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path"
	"time"
	pwsh "web-dashboard/powershell-utils"

	"github.com/labstack/echo/v4"
)

var editor = "C:/Program Files/Neovim/bin/nvim.exe"

func Ping(c echo.Context) error {
	return c.String(http.StatusOK, "OK")
}
func Edit(c echo.Context) error {
	start := time.Now()
	leaf := c.QueryParam("fileName") // Get-ChildItem "$env:CLOUD_DIR/Catalog/*.md"
	fp := path.Join(os.Getenv("CLOUD_DIR"), "Catalog", leaf)
	fmt.Println(editor, fp)
	//fp = "C:/Users/eshaa/draft.lua"
	open_cmd := fmt.Sprintf("wt.exe -d $env:CLOUD_DIR nvim %v", fp)
	pwsh.ExecPwshCmd(open_cmd)
	fmt.Println("[Debug] (", c.Path(), ") request elapsed time:", time.Since(start))
	return c.String(http.StatusOK, "success")
}

func RecentEagleItems(c echo.Context) error {
	//  TODO: IMPROVE SPEED (791.9206ms)
	a := pwsh.RunPwshCmd("./powershell-utils/recentEagleItems.ps1")
	return c.JSON(200, a)
}

func RecentNotes(c echo.Context) error {
	a := pwsh.RunPwshCmd("./powershell-utils/recentNotes.ps1")
	return c.JSON(200, a)
}

func NumTabs(c echo.Context) error {
	start := time.Now()
	a := pwsh.RunPwshCmd("./powershell-utils/waterfoxTabs.ps1")
	fmt.Println("[Debug] (", c.Path(), ") request elapsed time:", time.Since(start))
	return c.JSON(200, a)
}

// populate funcs
// set default params with pathparams
// they are empty strings not nil null case.
/*
Every PopulateFunction returns data that will be consumed by a template.
using the context, we can extract parameters or default arguments we can pass to the API calls.

if there is an error, the .error member is populated. This is checked first in the template and if it exists, the template is populated in the error case.

some of the populate functions bubbled errors by just returning c.string(400,err) which is less flexible.
*/

func PopulateGetNotesDetail(c echo.Context, templateName string) interface{} {
	return GetNotesNamesDates(-1, 0)
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
