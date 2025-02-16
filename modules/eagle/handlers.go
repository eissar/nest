package eaglemodule

import (
	"fmt"
	"github.com/eissar/nest/config"
	"github.com/eissar/nest/fileUtils"
	//"github.com/eissar/nest/modules/pwsh"

	"github.com/labstack/echo/v4"
)

func ServeThumbnailHandler(cfg *config.NestConfig) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := EagleItemId(c.Param("itemId"))
		if !id.IsValid() {
			res := fmt.Sprintf("get path=%s err=id of `%s` is not valid.", c.Path(), id)
			return c.JSON(200, res)
		}

		thumbnail, err := GetEagleThumbnail(cfg, id)
		if err != nil {
			res := fmt.Sprintf("get path=%s err=%s", c.Path(), err.Error())
			return c.String(200, res)
		}
		err = fileUtils.PathExists(thumbnail.ThumbnailPath)
		if err != nil {
			res := fmt.Sprintf("get path=%s err=%s", c.Path(), err.Error())
			return c.String(200, res)
		}
		// filepath exists.
		return c.File(thumbnail.ThumbnailPath)
		//return c.JSON(200, thumb)
	}
}
