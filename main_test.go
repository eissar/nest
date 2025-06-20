package main

import (
	"fmt"
	"testing"

	"github.com/eissar/nest/api"
	"github.com/eissar/nest/config"
)

func TestInfo(t *testing.T) {
	cfg := config.GetConfig()
	inf, err := api.ItemInfo(cfg.BaseURL(), "MBTR8K3VY1WR9")
	if err != nil {
		t.Error("testinfo", err.Error())
	}
	fmt.Print(inf)
}
