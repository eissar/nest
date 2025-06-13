package commandline

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"

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

func add() error {
	return nil
}

// TODO: remove via filter?
func CmdRemove() *cobra.Command {
	var removeItemIds []string
	// write
	removeCmd := &cobra.Command{
		Use:   "remove [ITEM_ID]",
		Short: "Moves selected Eagle items to the trash",

		Args: cobra.MinimumNArgs(1),

		PreRunE: func(cmd *cobra.Command, args []string) error {
			// Ensure all required arguments are set
			if len(args) == 0 {
				return errors.New("at least one item ID is required")
			}

			for _, arg := range args {
				if strings.TrimSpace(arg) == "" {
					return errors.New("item ID should not be empty")
				}
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := config.GetConfig()

			if err := removeItem(cfg.BaseURL(), removeItemIds); err != nil {
				return fmt.Errorf("error removing items: %w", err)
			}

			info, err := api.ItemInfo(cfg.BaseURL(), removeItemIds[0])
			fmt.Println(info)
			fmt.Printf("Items with IDs %v moved to Trash.\n", removeItemIds)
			if err != nil {
				return fmt.Errorf("error removing items: %w", err)
			}
			return nil
		},
		Aliases: []string{"rm"},
	}

	return removeCmd
}

// TODO: add flag to delete file after adding.
func CmdAdd() *cobra.Command {
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

			if err := api.ItemAddFromPath(cfg.BaseURL(), opts); err != nil {
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
func Adds() *cobra.Command {
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

// List creates the "list" command.
func List() *cobra.Command {
	// This variable will hold the value from the --limit flag.
	var limit int
	var filter string

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "Lists items from the Eagle library",
		Long:  `Retrieves and prints a list of items from the Eagle library in JSON format.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := config.GetConfig()
			opts := api.ItemListOptions{
				Limit:   limit,
				Keyword: filter,
			}

			list := func(cfg config.NestConfig, opts api.ItemListOptions) {
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
			list(cfg, opts)
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
	listCmd.Flags().StringVarP(&filter, "filter", "f", "", "Filter items by keyword(s)")

	return listCmd
}
func Reveal() *cobra.Command {
	revealCmd := &cobra.Command{
		Use:   "reveal [FILEPATH | ITEM_ID]",
		Short: "Reveals a file in the file explorer",
		Long: `Reveals a file in the system's file explorer.

You can provide a direct path to a file on your system.
Alternatively, you can provide an Eagle item ID, and the command
will resolve it to the item's location within the library.`,

		// This ensures exactly one argument is provided.
		Args: cobra.ExactArgs(1),

		RunE: func(cmd *cobra.Command, args []string) error {
			target := args[0]
			cfg := config.GetConfig()

			var err error

			resolveOrGetFilepath := func() (resolvedPath string) {
				resolvedPath, _ = filepath.Abs(target)
				if _, err := os.Stat(resolvedPath); err != nil {
					resolvedPath, err := nest.GetEagleThumbnailFullRes(&cfg, target)
					if err != nil {
						log.Fatalf("error getting thumbnail: %s", err.Error())
					}
					fmt.Printf("resolvedPath: %v\n", resolvedPath)
					return resolvedPath
				}

				return resolvedPath
			}

			err = launch.Reveal(resolveOrGetFilepath())
			return err
		},
	}

	return revealCmd
}
func Shutdown() *cobra.Command {
	cfg := config.GetConfig()
	shutdownCmd := &cobra.Command{
		Use:   "shutdown",
		Short: "Shuts down the running Eagle API server",
		Long:  `Sends a close request to the Eagle API server to shut it down gracefully.`,

		PreRunE: func(cmd *cobra.Command, args []string) error {

			pingEndpoint := fmt.Sprintf("http://localhost:%v/api/ping", cfg.Nest.Port)
			if !isServerRunning(pingEndpoint) {
				//not running
				return fmt.Errorf("shutdown: request to %s failed. The server is not running.\n", pingEndpoint)
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			closeEndpoint := fmt.Sprintf("http://localhost:%v/api/server/close", cfg.Nest.Port)

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

		},
	}
	return shutdownCmd
}
func Switch() *cobra.Command {
	cfg := config.GetConfig()
	switchCmd := &cobra.Command{
		Use:   "switch [LIBRARY_NAME]",
		Short: "Switches the active Eagle library",
		Long:  `Switches to a different Eagle library from your history by its name.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			libraryName := args[0]
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
			// TODO: should we instead print nothing?
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
					return nil
				}
				lib = strings.TrimSuffix(lib, ".LIBRARY")
				if libraryName == lib {
					switchTo(recentLibraries[i])
					return nil
				}
			}
			return nil
		},
	}
	return switchCmd
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

func CmdCobra() {
	var rootCmd = &cobra.Command{Use: "nest"}
	rootCmd.AddCommand(CmdAdd())
	rootCmd.AddCommand(Adds())
	rootCmd.AddCommand(List())
	rootCmd.AddCommand(Reveal())
	rootCmd.AddCommand(Switch())
	rootCmd.AddCommand(Shutdown())
	rootCmd.AddCommand(
		&cobra.Command{
			Use: "start",
			Run: func(cmd *cobra.Command, args []string) {
				core.Start() // blocking
			},
		})

	if err := rootCmd.Execute(); err != nil {
		// Cobra prints the error, so we just need to exit.
		os.Exit(1)
	}
}

// TODO: add to command line (all below)
func removeItem(baseUrl string, itemIds []string) error {
	err := api.ItemMoveToTrash(baseUrl, itemIds)
	if err != nil {
		return fmt.Errorf("failed to move item to trash: %s", err)
	}

	fmt.Println("Item moved to Trash successfully")
	return nil
}
