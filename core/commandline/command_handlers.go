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

	"github.com/eissar/nest/api"
	"github.com/eissar/nest/config"
	f "github.com/eissar/nest/format"
	"github.com/eissar/nest/plugins/launch"
	"github.com/eissar/nest/plugins/nest"
	"github.com/spf13/cobra"
)

// List creates the "list" command.
func List() *cobra.Command {
	// This variable will hold the value from the --limit flag.
	var limit int
	var filter string
	var properties string
	var o f.FormatType // output format

	// allowedFormats := []string{"json"}
	//allowedFormats := []string{"json", "log", "logfmt"}
	// TODO: url and website?
	defaultFields := []string{"id", "name", "tags", "annotation", "url", "website"}

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "Lists items from the Eagle library",
		Long:  "Retrieves and prints a list of items from the Eagle library in logfmt.\n",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := config.GetConfig()
			opts := api.ItemListOptions{
				Limit:   limit,
				Keyword: filter,
			}

			data, err := api.ItemList(cfg.BaseURL(), opts)
			if err != nil {
				return fmt.Errorf("while retrieving items: %w", err)
			}

			f.Format(o, data)
			return nil
		},
	}

	// TODO: improve argument parsing to accept commas like "1, 2" and "1,2"
	// NOTE: we won't be able to accept ', ' sep values without more advanced argument
	// handling. cobra uses <https://github.com/spf13/pflag> for parsing.
	// we could make properties a positional argument, but I don't see the benefit
	listCmd.Flags().IntVarP(&limit, "limit", "l", 10, "The maximum number of items to return")
	listCmd.Flags().StringVarP(&filter, "filter", "f", "", "Filter items by keyword(s)")
	listCmd.Flags().StringVarP(&properties, "properties", "p", "", "select properties to include in the output: "+f.HelpFmt(&api.ListItem{})+" default:"+f.HelpFmt(&defaultFields))

	listCmd.Flags().VarP(&o, "format", "o", "output format")

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
	var timeout int

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
				Timeout: time.Duration(timeout) * time.Second,
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

	shutdownCmd.Flags().IntVarP(&timeout, "timeout", "t", 10, "The maximum amount of time to wait for library to switch.")

	return shutdownCmd
}

// TODO: add flag to delete file after adding.
// TODO: rename -f file flag to -p path ?
// TODO: possible to check for item existence before api request?
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
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				// error if both a positional arg AND --file flag are used.
				if cmd.Flags().Changed("file") {
					return errors.New("cannot use both a positional argument and the --file flag")
				}
				addPath = args[0]
			}

			if addPath == "" {
				return errors.New("a filepath must be provided either as an argument or with the --file flag")
			}

			var err error
			addPath, err = filepath.Abs(addPath)
			if err != nil {
				return err
			}

			_, err = os.Stat(addPath)
			if err != nil {
				if errors.Is(err, os.ErrNotExist) {
					return fmt.Errorf("addPath: invalid path at %s (invalid or unavailable filepath)", addPath)
				}
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
				return err
			}

			fmt.Println("Item added successfully")
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
func Folder() *cobra.Command {
	var folderName string
	var folderOutput bool
	folderCmd := &cobra.Command{
		Use:   "folder [name]",
		Short: "Create a new folder",
		Long: `Create a new folder on the remote server.

You can specify the folder name either as a positional argument or using the --name flag.
If both are provided, the positional argument takes precedence.`,
		Example: `
  # Create a folder using a positional argument
  nest folder Reports

  # Create a folder using the --name flag
  nest folder --name Reports

  # Using shorthand for the flag
  nest folder -n Reports
`,
		Args: cobra.MaximumNArgs(1),

		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				folderName = args[0]
			} else if folderName == "" {
				return fmt.Errorf("create folder: flag 'name' cannot be nil and no positional param.")
			}
			cfg := config.GetConfig()
			if out, err := api.FolderCreate(cfg.BaseURL(), folderName); err != nil {
				return fmt.Errorf("create folder: %w", err)
			} else if folderOutput { // user wants folder id in output.
				fmt.Fprint(os.Stdout, out.ID)
			}
			return nil
		},
	}

	folderCmd.Flags().BoolVarP(&folderOutput, "output", "o", false, "whether or not to output folder ID on success.")
	folderCmd.Flags().StringVarP(&folderName, "name", "n", "", "Set a custom name for the folder")
	return folderCmd
}
func RecentLibraries() *cobra.Command {
	// var libraryList string
	librariesCmd := &cobra.Command{
		Use:   "libraries",
		Short: "list recent libraries",
		RunE: func(cmd *cobra.Command, args []string) error {
			// if len(args) > 0
			cfg := config.GetConfig()
			//
			// fmt.Print(cfg.Libraries)

			recentLibraries, err := api.LibraryHistory(cfg.BaseURL())
			if err != nil {
				log.Fatalf("could not retrieve recent libaries err=%s", err.Error())
			}

			if err != nil {
				return fmt.Errorf("recent libraries: %w", err)
			}

			var formattedLibs []string

			for _, v := range recentLibraries {
				// replace multiple [seperators] to a single one then conv to forward slashes `/`
				fmt := filepath.ToSlash(filepath.Clean(v))
				formattedLibs = append(formattedLibs, fmt)
			}

			// jsonFmtStdOut(cmd, recentLibraries, nil)

			stdout := json.NewEncoder(os.Stdout)

			stdout.Encode(formattedLibs)

			return nil // Added missing return statement
		},
	}

	return librariesCmd
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

// TODO: add nest list available or nest list list or nest list show or nest list libs or nest list libraries
// to show libraries available to switch to (print)
func Switch() *cobra.Command {
	var timeout int
	var targetLibraryName string

	cfg := config.GetConfig()

	// TODO: use instead cfg.Libraries
	recentLibraries, err := api.LibraryHistory(cfg.BaseURL())
	if err != nil {
		log.Fatalf("could not retrieve recent libaries err=%s", err.Error())
	}

	switchTo := func(libraryPath string) {
		err := nest.LibrarySwitchSync(cfg.BaseURL(), libraryPath, timeout)

		// err := api.SwitchLibrary(cfg.BaseURL(), libraryPath)
		if err != nil {

			if errors.Is(err, api.LibraryIsAlreadyTargetErr) {
				fmt.Println(err.Error())
			}

		}
	}
	switchCmd := &cobra.Command{
		Use:   "switch [LIBRARY_NAME]",
		Short: "Switches the active Eagle library",
		Long:  `Switches to a different Eagle library from your history by its name.`,
		Args:  cobra.ExactArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if !cmd.Flags().Changed("name") {
				// flag was not set with flag
				targetLibraryName = args[0]
			}

			if targetLibraryName == "" {
				return errors.New("name should not be an empty string")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			targetLibraryName = strings.ToUpper(targetLibraryName)
			for i, lib := range recentLibraries {
				lib = strings.ToUpper(lib)

				_, lib = filepath.Split(lib)
				if targetLibraryName == lib {
					switchTo(recentLibraries[i])
					return nil
				}
				lib = strings.TrimSuffix(lib, ".LIBRARY")
				if targetLibraryName == lib {
					libPath := recentLibraries[i]
					_, err := os.Stat(libPath)
					if err != nil {
						if errors.Is(err, os.ErrNotExist) {
							return fmt.Errorf("switch: invalid library path at %s (invalid or unavailable filepath)", libPath)
						}
					}

					switchTo(recentLibraries[i])
					return nil
				}
			}
			// not found in recentLibraries
			return errors.New("invalid library name: library must be present in recentLibraries (check available libraries in eagle GUI)")

			// return nil
		},
	}

	switchCmd.Flags().IntVarP(&timeout, "timeout", "t", 10, "The maximum amount of time to wait for library to switch.")
	switchCmd.Flags().StringVarP(&targetLibraryName, "name", "n", "", "name of the library to switch to")

	return switchCmd
}
