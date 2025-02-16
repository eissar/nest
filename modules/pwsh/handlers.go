package pwsh

import "github.com/labstack/echo/v4"

// echo.HandlerFunc
func RecentEagleItems(c echo.Context) error {
	//  TODO: IMPROVE SPEED (791.9206ms)
	a := RunScript("./recentEagleItems.ps1")
	return c.JSON(200, a)
}
