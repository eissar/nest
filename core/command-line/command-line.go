package commandline

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/spf13/cobra"
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

type addCmdOption struct {
	FlagName     string
	DefaultValue string
	Description  string
}

func Cmd() {

	addCmd := flag.NewFlagSet("add", flag.ExitOnError)

	// var addCmdOpts = []addCmdOption{
	// 	{FlagName: "file", DefaultValue: "", Description: "filepath that will be added to eagle"},
	// 	{FlagName: "name", DefaultValue: "", Description: "name"},
	// 	{FlagName: "website", DefaultValue: "", Description: "website"},
	// 	{FlagName: "annotation", DefaultValue: "", Description: "annotation"},
	// 	// { FlagName:    "tags", DefaultValue: "", Description: "tags", },
	// 	{FlagName: "folderid", DefaultValue: "", Description: "folderid"},
	// }

	addsCmd := flag.NewFlagSet("adds", flag.ExitOnError)
	// listCmd := flag.NewFlagSet("list", flag.ExitOnError)
	// listLimit := listCmd.Int("limit", 5, "number of items to retrieve")

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
		addPath := addCmd.String("file", "", "filepath that will be added to eagle")
		addName := addCmd.String("name", "", "name")
		addWebsite := addCmd.String("website", "", "website")
		addAnnotation := addCmd.String("annotation", "", "annotation")
		//addTags := addCmd.String("tags", "", "tags")
		addFolderId := addCmd.String("folderid", "", "folderid")

		addCmd.Parse(os.Args[2:])
		fmt.Println("amount positional: %s")

		// // numPositionalArgs := len(addCmd.Args())
		// for now, just assign first positional argument to path if -path not specified.
		if *addPath == "" {
			*addPath = addCmd.Arg(0)
		}

		opts := api.ItemAddFromPathOptions{Path: *addPath, Name: *addName, Website: *addWebsite, Annotation: *addAnnotation, FolderId: *addFolderId}

		err := Add1(cfg, opts)
		if err != nil {
			log.Fatalf("error adding item: %s", err.Error())
		}

	case "adds":
		cfg := config.GetConfig()
		addsCmd.Parse(os.Args[2:])
		var filepaths []string
		filepaths = addsCmd.Args()
		Adds(cfg, filepaths)
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

// TODO: add flag to delete file after adding.
func NewAdd() *cobra.Command {
	// These variables will hold the values from the flags.
	var addName, addWebsite, addAnnotation, addFolderId string
	var addPath string
	addCmd := &cobra.Command{
		Use:   "add [FILEPATH]",
		Short: "Adds a file to Eagle",
		Long: `Adds a file to your Eagle library with optional metadata.

The path to the file can be provided as the first argument directly
or by using the --file flag.`,

		Args: cobra.MaximumNArgs(1), // Allow zero or one positional argument. Error if more than one.
		// BEGIN
		PreRunE: func(cmd *cobra.Command, args []string) error {
			// Logic to determine the final path.
			// If a positional argument is given, it's the path.
			if len(args) > 0 {
				// Prevent confusion: error if both a positional arg AND --file flag are used.
				if cmd.Flags().Changed("file") {
					return errors.New("cannot use both a positional argument and the --file flag")
				}
				addPath = args[0]
			}

			// After checking flags and args, we must have a path.
			if addPath == "" {
				return errors.New("a filepath must be provided either as an argument or with the --file flag")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {

			cfg := config.GetConfig()

			opts := api.ItemAddFromPathOptions{
				Path:       addPath,
				Name:       addName,
				Website:    addWebsite,
				Annotation: addAnnotation,
				FolderId:   addFolderId,
			}

			// Your business logic function
			if err := Add1(cfg, opts); err != nil {
				// Use fmt.Errorf to wrap the error for more context
				return fmt.Errorf("error adding item: %w", err)
			}

			fmt.Println("Item added successfully!")
			return nil
		},
	}
	// Define all the flags for the 'add' command
	addCmd.Flags().StringVarP(&addPath, "file", "f", "", "Filepath to add to Eagle")
	addCmd.Flags().StringVarP(&addName, "name", "n", "", "Set a custom name for the item")
	addCmd.Flags().StringVar(&addWebsite, "website", "", "Set a source website URL")
	addCmd.Flags().StringVar(&addAnnotation, "annotation", "", "Add an annotation or description")
	addCmd.Flags().StringVar(&addFolderId, "folderid", "", "ID of the folder to add the item into")

	return addCmd

}
func NewAdds() *cobra.Command {
	addsCmd := &cobra.Command{
		Use:   "adds [FILE1] [FILE2]",
		Short: "Adds multiple files to Eagle in a single batch",
		Long: `Adds one or more files to your Eagle library without metadata.
Provide the paths to the files as arguments separated by spaces.`,

		// This validator ensures at least one positional argument is given.
		// It automatically provides a user-friendly error message if the condition fails.
		Args: cobra.MinimumNArgs(1),

		RunE: func(cmd *cobra.Command, args []string) error {
			// 'args' is now a slice containing all the file paths provided.
			pths := args

			cfg := config.GetConfig()
			opts := []api.ItemAddFromPathOptions{}

			for _, v := range pths {
				opts = append(opts, api.ItemAddFromPathOptions{Path: v})
			}

			if err := api.ItemAddFromPaths(cfg.BaseURL(), opts); err != nil {
				return fmt.Errorf("error while adding items in batch: %w", err)
			}

			fmt.Printf("Successfully processed %d items.\n", len(pths))
			return nil
		},
	}

	return addsCmd
}

// NewList creates the "list" command.
func NewList() *cobra.Command {
	// This variable will hold the value from the --limit flag.
	var limit int

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "Lists items from the Eagle library",
		Long:  `Retrieves and prints a list of items from the Eagle library in JSON format.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := config.GetConfig()
			opts := api.ItemListOptions{
				Limit: limit,
			}

			List(cfg, opts)
			data, err := api.ItemList(cfg.BaseURL(), opts)
			if err != nil {
				return fmt.Errorf("while retrieving items: %w", err)
			}

			// output, err := json.MarshalIndent(data, "", "  ") // Using MarshalIndent for nice formatting
			// if err != nil {
			// 	return fmt.Errorf("while parsing list items: %w", err)
			// }

			// On success, print the JSON to standard output.
			fmt.Println(data)
			return nil
		},
	}

	// Define the --limit flag, with a short version -l, a default value, and a help message.
	listCmd.Flags().IntVarP(&limit, "limit", "l", 10, "The maximum number of items to return")

	return listCmd
}

func CmdCobra() {
	var rootCmd = &cobra.Command{Use: "nest"}
	rootCmd.AddCommand(NewAdd())
	rootCmd.AddCommand(NewAdds())
	rootCmd.AddCommand(NewList())

	if err := rootCmd.Execute(); err != nil {
		// Cobra prints the error, so we just need to exit.
		os.Exit(1)
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

func Add1(cfg config.NestConfig, item api.ItemAddFromPathOptions) error {
	fmt.Println("adding...")
	err := api.ItemAddFromPath(cfg.BaseURL(), item)
	return err
}

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

func List(cfg config.NestConfig, opts api.ItemListOptions) {
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
