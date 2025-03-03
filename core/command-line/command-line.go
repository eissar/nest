package commandline

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/eissar/nest/api"
	"github.com/eissar/nest/config"
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

// validateIsEagleServerRunning checks if the Eagle server is running at the specified URL.
func isServerRunning(url string) bool {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return false //, fmt.Errorf("error creating request: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return false //, fmt.Errorf("error making request: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return false //, fmt.Errorf("received code other than 200: %v", resp.StatusCode)
	}

	return true //, nil
}

// returns resp or calls log.fatal
func Shutdown(cfg config.NestConfig) {
	closeEndpoint := fmt.Sprintf("http://localhost:%v/api/server/close", cfg.Nest.Port)
	pingEndpoint := fmt.Sprintf("http://localhost:%v/api/ping", cfg.Nest.Port)
	if !isServerRunning(pingEndpoint) {
		//not running
		log.Fatalf("shutdown: request to %s failed. The server is not running.", pingEndpoint)
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest("GET", closeEndpoint, nil)
	if err != nil {
		log.Fatalf("error creating request: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("error making request: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("received code other than 200: %v", resp.StatusCode)
	}

}
