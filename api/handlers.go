package api

import (
	"github.com/labstack/echo/v4"
)

func handleAddItemFromUrl(c echo.Context) error {
	a := c.QueryParams()
	if a.Has("url") && a.Has("name") {
		return c.String(200, a.Get("url"))
	}
	// url Required，the URL of the image to be added. Supports http、 https、 base64
	// name Required，The name of the image to be added.
	// website The Address of the source of the image
	// tags Tags for the image.
	// star The rating for the image.
	// annotation The annotation for the image.
	// modificationTime The creation date of the image. The parameter can be used to alter the image's sorting order in Eagle.
	// folderId If this parameter is defined, the image will be added to the corresponding folder.
	// headers Optional, customize the HTTP headers properties, this could be used to circumvent the security of certain websites.

	//data, err := AddItemFromURL(url)
	return c.String(200, "missing mandatory parameters")
}
