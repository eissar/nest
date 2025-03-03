package browser

// go get -u github.com/eissar/browser-query@master
import (
	"sync"

	"github.com/eissar/browser-query"
)

/*
	type uploadTabsBody struct {
		Body string `json:"body"`
	}

	func UploadTabs(c echo.Context) error {
		body := c.Request().Body

		var bytes []byte
		if _, err := body.Read(bytes); err != nil {
			msg := fmt.Sprintf("uploadtabs: could not read request body err=%s", err.Error())
			return c.String(400, msg)
		}

		//fmt.Println("[SUCCESS]", c)
		return c.String(200, "OK")
	}
*/

// TODO: ? merge from github.com/eissar/browser-query
// WARN: browserquery has state, bad library; reformat before using here.

var (
	clientsMu sync.RWMutex
	clients   = make(map[*browserQuery.Client]bool)
)

func t() {
	//main.HandleSSE
}
