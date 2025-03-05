package nest

import (
	"fmt"
	"testing"

	"github.com/eissar/nest/config"
)

func TestThumb(t *testing.T) {
	cfg := config.GetConfig()
	i, err := GetEagleThumbnailFullRes(&cfg, "LVKFORWVA0O4O")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(i)
}
func TestList(t *testing.T) {
	cfg := config.GetConfig()
	i, err := GetList(&cfg)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(i)
}
func TestCurrentLibrary(t *testing.T) {
	//cfg := config.GetConfig()
	libPath, err := CurrentLibraryPath()
	if err != nil {
		t.Fatalf("%s", err.Error())
	}
	fmt.Println(libPath)
}
func TestSwitchLibrary(t *testing.T) {
	cfg := config.GetConfig()
	err := LibrarySwitchSync(cfg.BaseURL(), cfg.Libraries.Paths[1])
	if err != nil {
		t.Fatalf("%s", err.Error())
	}
}
