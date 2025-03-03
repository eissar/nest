package nest

import (
	"fmt"

	"github.com/eissar/nest/config"

	//"github.com/eissar/nest/plugins/pwsh"

	"github.com/labstack/echo/v4"
)

func ServeThumbnailHandler(cfg *config.NestConfig) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("itemId")
		resFlag := c.QueryParam("fq")

		getThumbnail := func() (string, error) {
			if resFlag == "true" {
				return GetEagleThumbnailFullRes(cfg, id)
			}
			return GetEagleThumbnail(cfg, id)
		}

		thumbnail, err := getThumbnail()
		if err != nil {
			res := fmt.Sprintf("get thumbnail path=%s err=%s", c.Path(), err.Error())
			return c.String(400, res)
		}
		// filepath exists.
		return c.File(thumbnail)
	}
}
