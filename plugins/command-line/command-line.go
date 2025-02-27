package commandline

import (
	"fmt"
	"log"
	"os"

	"github.com/eissar/nest/config"
	"github.com/eissar/nest/eagle/api"
)

func Add(cfg config.NestConfig, pth *string) {
	if len(*pth) == 0 {
		log.Fatalf("[ERROR] add: flag `-file` is required.")
	}
	fmt.Println("path:", *pth)

	err := api.AddItemFromPath(cfg.FmtURL(), *pth)
	if err != nil {
		log.Fatalf("[ERROR] while adding eagle item: err=%s", err.Error())
	}
	os.Exit(0)
}
