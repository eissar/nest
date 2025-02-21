package api // Same package as the source file

import (
	"fmt"
	"strconv"
	"testing"
	//"time"

	"github.com/eissar/nest/config"
)

/*
:tcd %:p:h
:!go test -run TestList
:!gotestsum --format testname -- -run ^TestList$
*/

var cfg = config.GetConfig()

// lists 1
func TestList(t *testing.T) {
	baseUrl := "http://" + cfg.Host + ":" + strconv.Itoa(cfg.Port)

	result, err := List(baseUrl, 1)
	if err != nil {
		t.Error(err)
	}

	fmt.Println(result)
}

// get count of all items in library

func TestListLengths(t *testing.T) {
	baseUrl := "http://" + cfg.Host + ":" + strconv.Itoa(cfg.Port)

	lens := []int{1, 5, 20, 200}

	for _, limit := range lens {
		result, err := List(baseUrl, limit)
		if err != nil {
			t.Fatalf("%s", err.Error())
		}
		l := len(result.Data)
		if limit != l {
			t.Fatalf("expected data of len %v, but got %v", limit, l)
		}
	}
}
