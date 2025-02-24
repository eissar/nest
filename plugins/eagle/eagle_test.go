package eaglemodule

import (
	"fmt"
	"testing"

	"github.com/eissar/nest/config"
)

func TestThumb(t *testing.T) {
	cfg := config.GetConfig()
	i, err := GetEagleThumbnailV2(&cfg, "LVKFORWVA0O4O")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(i)
}
func TestList(t *testing.T) {
	cfg := config.GetConfig()
	i, err := GetListV0(&cfg)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(i)

}
