package browsermodule

import (
	"database/sql"
	"fmt"
	_ "fmt"
	"os"
	_ "strings"
	"time"

	"github.com/labstack/echo/v4"
	_ "github.com/mattn/go-sqlite3"
)

func RegisterRoutesFromGroup(g *echo.Group) {
	g.GET("/test", func(c echo.Context) error {
		prof := os.Getenv("firefox_profile_1")
		if prof == "" {
			fmt.Println("[Error] env variable firefox_profile_1 cannot be empty or nonexistent.")
			return nil
		}
		conn_str := fmt.Sprintf("%s/places.sqlite?access_mode=read_only", prof)
		fmt.Println(conn_str)
		db, err := sql.Open("sqlite3", conn_str)
		defer db.Close()

		if err != nil {
			return c.String(400, err.Error())
		}

		start := time.Now()
		// ./TopOriginsByFrecency.sql
		resp, err := db.Query("SELECT * FROM moz_origins " +
			"WHERE frecency IS NOT NULL " +
			"ORDER BY frecency DESC " +
			"LIMIT 10;")
		if err != nil {
			return c.String(400, err.Error())
		}
		// fmt.Printf("resp: %v\n", resp)

		// TODO: better error handling
		//var responseColumns string
		for {
			if resp.Next() != true {
				break
			}
			var id string
			var host string
			var frecency int
			// ignore cols we don't want.
			var nill *string // nill can be nil
			//id,prefix,host,frecency,recalc_frecency,alt_frecency,recalc_alt_frecency
			err = resp.Scan(&id, &nill, &host, &frecency, &nill, &nill, &nill)

			if err != nil {
				fmt.Println(err.Error())
				fmt.Println(id)
			}
			fmt.Println("[SUCCESS]", host, frecency)

		}
		fmt.Println("[LOG] elapsed:", time.Since(start))

		return c.String(200, "OK")
	})
}

// func GetBrowserDomains

// func GetBrowserTabsCount
