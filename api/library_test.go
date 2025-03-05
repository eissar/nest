package api

import (
	"fmt"
	"testing"

	"github.com/eissar/nest/config"
)

func TestSwitchLibrary(t *testing.T) {
	cfg := config.GetConfig()
	lib := cfg.Libraries.Paths[1]

	err := SwitchLibrary(cfg.BaseURL(), lib)
	if err != nil {
		t.Fatalf("couldn't switch lib=%s err=%v", lib, err)
	}
}

// demonstrates that running functions immediately after
// switching will give you stale data -unsafe!
func TestAsyncSwitchLibrary(t *testing.T) {
	cfg := config.GetConfig()
	lib := cfg.Libraries.Paths[1]

	err := SwitchLibrary(cfg.BaseURL(), lib)
	if err != nil {
		t.Fatalf("couldn't switch lib=%s err=%v", lib, err)
	}
	_, err = GetLibraryInfo(cfg.BaseURL())
	if err != nil {
		t.Fatalf("couldn't get lib info err=%v", err)
	}
}

/*
func TestSyncSwitchLibrary(t *testing.T) {
	cfg := config.GetConfig()
	lib := cfg.Libraries.Paths[0]

	err := SwitchLibrarySync(cfg.BaseURL(), lib)
	if err != nil {
		t.Fatalf("couldn't switch lib=%s err=%v", lib, err)
	}
}
*/

func TestGetLibraryInfo(t *testing.T) {
	cfg := config.GetConfig()

	_, err := GetLibraryInfo(cfg.BaseURL())
	if err != nil {
		t.Fatalf("couldn't get lib info err=%v", err)
	}
}

func TestGetRecent(t *testing.T) {
	cfg := config.GetConfig()

	libs, err := Recent(cfg.BaseURL())
	if err != nil {
		t.Fatalf("couldn't get recent libraries err=%v", err)
	}

	fmt.Println(libs)
}
