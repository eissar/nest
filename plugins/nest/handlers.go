package nest

import (
	"fmt"
	"os"

	"github.com/eissar/nest/config"
	"github.com/eissar/nest/eagle/api"

	//"github.com/eissar/nest/plugins/pwsh"

	"github.com/labstack/echo/v4"
)

func ServeThumbnailHandler(cfg *config.NestConfig) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("itemId")
		if !api.IsValidItemID(id) {
			res := fmt.Sprintf("get path=%s err=id of `%s` is not valid.", c.Path(), id)
			return c.JSON(200, res)
		}

		thumbnail, err := GetEagleThumbnail(cfg, id)
		if err != nil {
			res := fmt.Sprintf("get path=%s err=%s", c.Path(), err.Error())
			return c.String(200, res)
		}
		_, err = os.Stat(thumbnail)
		if err != nil {
			res := fmt.Sprintf("get path=%s err=%s", c.Path(), err.Error())
			return c.String(200, res)
		}
		// filepath exists.
		return c.File(thumbnail)
	}
}
