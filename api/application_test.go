package api

import (
	"fmt"
	"testing"

	"github.com/eissar/nest/config"
)

// error check
func TestApplicationInfo(t *testing.T) {
	cfg := config.GetConfig()
	resp, err := ApplicationInfo(cfg.BaseURL())
	if err != nil {
		t.Fatalf("could not get application info err=%s", err.Error())
	}

	fmt.Println(resp)
}
