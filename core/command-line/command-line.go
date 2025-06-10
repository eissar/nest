package commandline

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/eissar/nest/api"
	"github.com/eissar/nest/config"
	"github.com/eissar/nest/core"
	"github.com/eissar/nest/plugins/launch"
	"github.com/eissar/nest/plugins/nest"
)

// show a simpler error if we know the error.
// we can do some redundant over-checking here and
// it shouldn't matter much.
func catchKnownErrors(err error) {
	if errors.Is(err, api.EagleNotOpenErr) {
		fmt.Println(api.EagleNotOpenErr.Error())
		os.Exit(1)
	}
}

func Cmd() {
	addCmd := flag.NewFlagSet("add", flag.ExitOnError)
	addPath := addCmd.String("file", "", "filepath that will be added to eagle")
	addName := addCmd.String("name", "", "name")
	// addWebsite := addCmd.String("website", "", "website")
	addAnnotation := addCmd.String("annotation", "", "annotation")
	//addTags := addCmd.String("tags", "", "tags")
	addFolderId := addCmd.String("folderid", "", "folderid")

	addsCmd := flag.NewFlagSet("adds", flag.ExitOnError)

	listCmd := flag.NewFlagSet("list", flag.ExitOnError)
	listLimit := listCmd.Int("limit", 5, "number of items to retrieve")

	revealCmd := flag.NewFlagSet("reveal", flag.ExitOnError)
	revealPath := revealCmd.String("target", "", "filepath or item id to reveal")

	switchCmd := flag.NewFlagSet("switch", flag.ExitOnError)
	switchName := switchCmd.String("name", "", "name of library to switch to.")
	//revealCmd := flag.NewFlagSet("reveal", flag.ExitOnError)

	help := flag.Bool("help", false, "print help information")
	start := flag.Bool("start", false, "run the utility server")
	//debug := flag.Bool("debug", true, "shows additional information in the console while running. (does nothing)")
	stop := flag.Bool("stop", false, "stop the utility server")
	flag.Parse()

	if *help || len(os.Args) < 2 {
		fmt.Println("expected flag or subcommand.")
		flag.PrintDefaults()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "add":
		cfg := config.GetConfig()
		addWebsite := addCmd.String("website", "", "website")
		addCmd.Parse(os.Args[2:])

		opts := api.ItemAddFromPathOptions{Path: *addPath, Name: *addName, Website: *addWebsite, Annotation: *addAnnotation, FolderId: *addFolderId}

		Add1(cfg, opts)

	case "adds":
		cfg := config.GetConfig()
		addsCmd.Parse(os.Args[2:])
		var filepaths []string
		filepaths = addsCmd.Args()
		Adds(cfg, filepaths)
		os.Exit(0)

	case "list":
		cfg := config.GetConfig()

		listCmd.Parse(os.Args[2:])
		List(cfg, listLimit)
		os.Exit(0)
	case "reveal":
		cfg := config.GetConfig()

		revealCmd.Parse(os.Args[2:])
		Reveal(cfg, revealPath)
		os.Exit(0)
	case "switch":
		cfg := config.GetConfig()
		switchCmd.Parse(os.Args[2:])

		if *switchName != "" {
			Switch(cfg, *switchName)
			os.Exit(0)
		}
		if len(os.Args) < 3 {
			log.Fatalf("must pass flag -name")
			flag.PrintDefaults()
			os.Exit(1)
		}
		Switch(cfg, os.Args[2])
		os.Exit(0)
	}

	//if *debug { } /* pwsh.ExecPwshCmd("./powershell-utils/openUrl.ps1 -Uri 'http://localhost:1323/app/notes'") */

	if *start {
		core.Start() //blocking
	}

	if *stop {
		err := Shutdown(config.GetConfig())
		if err != nil {
			fmt.Printf("stop: %s", err.Error())
		}
		os.Exit(0)
	}
}

func Add(cfg config.NestConfig, pth *string) {
	if pth == nil || *pth == "" {
		log.Fatalf("[ERROR] add: flag `-file` is required.")
	}

	opts := api.ItemAddFromPathOptions{Path: *pth}

	err := api.ItemAddFromPath(cfg.BaseURL(), opts)
	if err != nil {
		log.Fatalf("Error while adding eagle item: err=%s", err.Error())
	}
}

func Add1(cfg config.NestConfig, item api.ItemAddFromPathOptions) {
	err := api.ItemAddFromPath(cfg.BaseURL(), item)
	if err != nil {
		log.Fatalf("Error while adding eagle item: err=%s", err.Error())
	}
}

// TODO: merge this stupid with other add
// TODO: Make processing continue after error; then report errors.
func Adds(cfg config.NestConfig, pths []string) {
	if len(pths) == 0 {
		log.Fatalf("[ERROR] adds: flag `-files` is required.")
	}

	opts := []api.ItemAddFromPathOptions{}

	for _, v := range pths {
		opts = append(opts, api.ItemAddFromPathOptions{Path: v})

	}
	err := api.ItemAddFromPaths(cfg.BaseURL(), opts)
	if err != nil {
		log.Fatalf("Error while adding eagle item: err=%s", err.Error())
	}

}

func List(cfg config.NestConfig, limit *int) {
	opts := api.ItemListOptions{
		Limit: *limit,
	}

	data, err := api.ItemList(cfg.BaseURL(), opts)
	if err != nil {
		log.Fatalf("[ERROR] list: while retrieving items: err=%s", err.Error())
	}

	output, err := json.Marshal(data)
	if err != nil {
		log.Fatalf("[ERROR] list: while parsing list items: err=%s", err.Error())
	}

	fmt.Fprintf(os.Stdout, "%v", string(output))
}

// param t string: target filepath or item id to reveal (in explorer)
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

// validateIsEagleServerRunning checks if the nest server is running at the specified URL.
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
func Shutdown(cfg config.NestConfig) error {
	closeEndpoint := fmt.Sprintf("http://localhost:%v/api/server/close", cfg.Nest.Port)
	pingEndpoint := fmt.Sprintf("http://localhost:%v/api/ping", cfg.Nest.Port)
	if !isServerRunning(pingEndpoint) {
		//not running
		return fmt.Errorf("shutdown: request to %s failed. The server is not running.\n", pingEndpoint)
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest("GET", closeEndpoint, nil)
	if err != nil {
		log.Fatalf("error creating request: %v", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("error making request: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("received code other than 200: %v", resp.StatusCode)
	}

	return nil
}

func Switch(cfg config.NestConfig, libraryName string) {
	if libraryName == "" {
		log.Fatalf("library name cannot be empty")
	}
	switchTo := func(libraryPath string) {
		err := nest.LibrarySwitchSync(cfg.BaseURL(), libraryPath)
		// err := api.SwitchLibrary(cfg.BaseURL(), libraryPath)
		if err != nil {
			log.Fatalf("could not switch library err=%s", err.Error())
		}
	}

	currLib, err := nest.CurrentLibrary()
	if err != nil {
		catchKnownErrors(err)
		log.Fatalf("unknown error getting current library err=%s", err.Error())
	}

	if currLib.Name == libraryName {
		log.Fatalf("library is already %s", libraryName)
	}

	recentLibraries, err := api.LibraryHistory(cfg.BaseURL())
	if err != nil {
		log.Fatalf("could not retrieve recent libaries err=%s", err.Error())
	}

	libraryName = strings.ToUpper(libraryName)
	for i, lib := range recentLibraries {
		lib = strings.ToUpper(lib)

		_, lib = filepath.Split(lib)
		if libraryName == lib {
			switchTo(recentLibraries[i])
			return
		}
		lib = strings.TrimSuffix(lib, ".LIBRARY")
		if libraryName == lib {
			switchTo(recentLibraries[i])
			return
		}
	}

}
