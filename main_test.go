package main

import (
	"fmt"
	"testing"

	"github.com/eissar/eagle-go"
	"github.com/eissar/nest/config"
)

func TestInfo(t *testing.T) {
	cfg := config.GetConfig()
	inf, err := eagle.ItemInfo(cfg.BaseURL(), "MBTR8K3VY1WR9")
	if err != nil {
		t.Error("testinfo", err.Error())
	}
	fmt.Print(inf)
}
