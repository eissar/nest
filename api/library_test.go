package api

import (
	"fmt"
	"testing"

	"github.com/eissar/nest/config"
)

//- [X] /api/library/info
//- [X] /api/library/history
//- [X] /api/library/switch
//- [X] /api/library/icon

// sanity test
func TestLibraryInfo(t *testing.T) {
	cfg := config.GetConfig()

	_, err := LibraryInfo(cfg.BaseURL())
	if err != nil {
		t.Fatalf("couldn't get lib info err=%v", err)
	}
}

// sanity test
func TestLibraryHistory(t *testing.T) {
	cfg := config.GetConfig()

	_, err := LibraryHistory(cfg.BaseURL())
	if err != nil {
		t.Fatalf("couldn't get recent libraries err=%v", err)
	}

	// fmt.Println(libs)
}

// sanity test
func TestLibrarySwitch(t *testing.T) {
	cfg := config.GetConfig()

	// get current library
	currentLibrary, err := LibraryInfo(cfg.BaseURL())
	if err != nil {
		t.Fatalf("error getting library info err=%s", err.Error())
	}

	// find a library that is not this library...
	targetLibrary := (func() string {
		for _, lib := range cfg.Libraries.Paths {
			if lib != currentLibrary.Library.Path {
				return lib
			}
		}
		t.Fatalf("couldn't switch libraries. no other target libraries found.")
		return ""
	}())

	// switch to that library
	err = LibrarySwitch(cfg.BaseURL(), targetLibrary)
	if err != nil {
		t.Fatalf("couldn't switch lib=%s err=%v", targetLibrary, err)
	}
}

// demonstrates that running functions immediately after
// switching will give you stale data - unsafe!
func TestLibrarySwitchAsync(t *testing.T) {
	cfg := config.GetConfig()
	lib := cfg.Libraries.Paths[1]

	err := LibrarySwitch(cfg.BaseURL(), lib)
	if err != nil {
		t.Fatalf("couldn't switch lib=%s err=%v", lib, err)
	}
	_, err = LibraryInfo(cfg.BaseURL())
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
