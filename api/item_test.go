package api // Same package as the source file

import (
	"fmt"
	"strconv"
	"testing"

	"net/http"
	"net/http/httptest"
	"net/url"

	//"time"

	"github.com/eissar/nest/config"
	"github.com/labstack/echo/v4"
)

/*
:tcd %:p:h
:!go test -run TestList
:!gotestsum --format testname -- -run ^TestList$
*/

var cfg = config.GetConfig()
var host = cfg.Host + ":" + strconv.Itoa(cfg.Port)

// lists 1
func TestList(t *testing.T) {
	baseUrl := "http://" + cfg.Host + ":" + strconv.Itoa(cfg.Port)

	result, err := ItemList(baseUrl, ItemListOptions{Limit: 0})
	if err != nil {
		t.Fatal(err.Error())
	}

	if len(result) != 1 {
		t.Fatalf("expected len 1 with parameter limit 0, instead got %v", len(result))
	}
	// fmt.Println(result)
}

// get count of all items in library

func TestListLengths(t *testing.T) {
	baseUrl := cfg.FmtURL()

	lens := []int{1, 5, 20, 200}

	for _, limit := range lens {
		result, err := ItemList(baseUrl, ItemListOptions{Limit: limit})
		if err != nil {
			t.Fatalf("%s", err.Error())
		}
		l := len(result)
		if limit != l {
			t.Fatalf("expected data of len %v, but got %v", limit, l)
		}
	}
}

func TestListWrapper(t *testing.T) {
	ep := "http://" + host + "/api/item/list"

	urls := []string{
		ep + "?limit=10",
	}

	for _, u := range urls {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, u, nil)
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)
		err := wrapperHandler(c)
		if err != nil {
			t.Fatalf("%s", err.Error())
		}
	}
}

func TestListInvoke(t *testing.T) {
	ep := "http://" + host + "/api/item/list"

	urls := []string{
		ep + "?limit=10",
	}

	for _, u := range urls {
		var err error
		req := httptest.NewRequest(http.MethodGet, u, nil)
		t.Fatalf("test")
		req.URL, err = url.Parse(u)
		//if err != nil {
		//	t.Fatalf("%s", err.Error())
		//}

		result, err := InvokeEagleAPIV1(req)
		if err != nil {
			t.Fatalf("%s", err.Error())
		}

		fmt.Println(result)
	}
}

func TestAddItemFromPath(t *testing.T) {
	baseUrl := cfg.FmtURL()
	err := AddItemFromPath(baseUrl, itemFromPath{Path: `C:/Users/eshaa/Downloads/download.jpg`})
	if err != nil {
		t.Fatal(err.Error())
	}
}

func TestItemInfo(t *testing.T) {
	x, err := ItemInfo(cfg.BaseURL(), "M7YCMVLJ090PF")
	if err != nil {
		t.Fatalf("%s", err)
	}
	fmt.Println(x)
}
