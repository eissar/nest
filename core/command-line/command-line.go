package commandline

import (
	"encoding/json"
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

	err := api.AddItemFromPath(cfg.BaseURL(), *pth)
	if err != nil {
		log.Fatalf("[ERROR] while adding eagle item: err=%s", err.Error())
	}
	os.Exit(0)
}

func List(cfg config.NestConfig, limit *int) {
	data, err := api.ListV2(cfg.BaseURL(), *limit)
	if err != nil {
		log.Fatalf("[ERROR] list: while retrieving items: err=%s", err.Error())
	}

	output, err := json.Marshal(data)
	if err != nil {
		log.Fatalf("[ERROR] list: while parsing list items: err=%s", err.Error())
	}

	fmt.Fprintf(os.Stdout, "%v", string(output))
}
