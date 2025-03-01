package commandline

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"

	"github.com/eissar/nest/config"
	"github.com/eissar/nest/eagle/api"
	"github.com/eissar/nest/plugins/launch"
	"github.com/eissar/nest/plugins/nest"
)

func Add(cfg config.NestConfig, pth *string) {
	if len(*pth) == 0 {
		log.Fatalf("[ERROR] add: flag `-file` is required.")
	}
	fmt.Println("path:", *pth)

	obj, err := api.ConstructItemFromPath(
		*pth,
	)
	fmt.Println("path:", obj.Path)
	if err != nil {
		log.Fatalf("[ERROR] while constructing request: err=%s", err.Error())
	}

	err = api.AddItemFromPath(cfg.BaseURL(), obj)
	if err != nil {
		log.Fatalf("[ERROR] while adding eagle item: err=%s", err.Error())
	}
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

// param t string: target filepath or item id to reveal
func Reveal(cfg config.NestConfig, t *string) {
	if len(*t) == 0 {
		log.Fatalf("[ERROR] add: flag `-target` is required.")
	}
	//fmt.Println("path:", *t)

	resolveOrGetFilepath := func() (resolvedPath string) {
		resolvedPath, _ = filepath.Abs(*t)
		if _, err := os.Stat(resolvedPath); err != nil {
			resolvedPath, err := nest.GetEagleThumbnailFullRes(&cfg, *t)
			if err != nil {
				log.Fatalf("error getting thumbnail: %s", err.Error())
			}
			resolvedPath, err = url.PathUnescape(resolvedPath)
			if err != nil {
				log.Fatalf("error cleaning thumbnail path: %s", err.Error())
			}
			fmt.Printf("resolvedPath: %v\n", resolvedPath)
			return resolvedPath
		}

		return resolvedPath
	}

	err := launch.Reveal(resolveOrGetFilepath())
	if err != nil {
		log.Fatalf("[ERROR] while adding eagle item: err=%s", err.Error())
	}
}
