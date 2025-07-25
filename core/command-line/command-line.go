package commandline

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"strings"
	"time"

	"github.com/go-logfmt/logfmt"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

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
	if errors.Is(err, api.EagleNotOpenOrUnavailableErr) {
		fmt.Println(api.EagleNotOpenOrUnavailableErr.Error())
		os.Exit(1)
	}
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

func logFmtStdOut(data []*api.ListItem, props []string) {
	outp := logfmt.NewEncoder(os.Stdout)

	outp.EncodeKeyvals("test", nil)
	outp.EndRecord()
	// for _, item := range data {
	//	// outp.EncodeKeyval("url", item.URL)
	//	// outp.EncodeKeyval("folderIds", strings.Join(item.Folders, ", "))
	//	outp.EndRecord()
	// }

}

// exclude properties to exclude
func jsonFmtStdOut(cmd *cobra.Command, data []*api.ListItem, exclude []string) error {
	json.NewEncoder(os.Stdout)

	b, err := json.Marshal(data)
	if err != nil {
		return err
	}

	var m []map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		return err
	}

	for i, _ := range m {
		for _, key := range exclude {
			delete(m[i], key)
		}
	}

	st := json.NewEncoder(os.Stdout)

	st.Encode(m)

	return nil
}

func rebuildCmdText(cmd *cobra.Command, args []string) string {
	var commandBuilder strings.Builder

	commandBuilder.WriteString(cmd.CommandPath())

	cmd.Flags().Visit(func(flag *pflag.Flag) {
		commandBuilder.WriteString(fmt.Sprintf(" --%s=%s", flag.Name, flag.Value))
	})
	if len(args) > 0 {
		commandBuilder.WriteString(" ")
		commandBuilder.WriteString(strings.Join(args, " "))
	}

	return commandBuilder.String()
}

// List creates the "list" command.
func List() *cobra.Command {
	// This variable will hold the value from the --limit flag.
	var limit int
	var filter string
	var properties string
	var format string

	allowedFormats := []string{"json"}
	//allowedFormats := []string{"json", "log", "logfmt"}
	// TODO: url and website?
	defaultFields := []string{"id", "name", "tags", "annotation", "url", "website"}

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "Lists items from the Eagle library",
		Long:  "Retrieves and prints a list of items from the Eagle library in logfmt.\n",

		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if slices.Contains(allowedFormats, format) {
				return nil
			}
			return fmt.Errorf("invalid format: %s. Please use one of: %s", format, helpFmt(&allowedFormats))
		},
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

			fmtFields := func(props string) []string {
				// strings.ReplaceAll(props, " ", "") // remove whitespace
				fields := strings.Split(props, ",")

				return fields
			}

			switch format {
			case "json":
				var targetProperties []string
				if properties != "" {
					targetProperties = fmtFields(properties)
				} else {
					targetProperties = defaultFields
				}

				allFields := structToKeys(&api.ListItem{}) // no struct keys should have any whitespace
				// find fields which are in allfields but not in inputFilterFields
				exclFields := filterFieldsByReference(targetProperties, allFields)
				err = jsonFmtStdOut(cmd, data, exclFields)
				if err != nil {
					fmt.Printf("jsonFmtStdOut: %v\n", err)
				}
			case "logfmt":
				logFmtStdOut(data, strings.Split(properties, ","))
			}

			// On success, print the JSON to standard output.
			// fmt.Println(output)
			return nil
		},
	}

	// TODO: improve argument parsing to accept commas like "1, 2" and "1,2"
	// NOTE: we won't be able to accept ', ' sep values without more advanced argument
	// handling. cobra uses <https://github.com/spf13/pflag> for parsing.
	// we could make properties a positional argument, but I don't see the benefit
	listCmd.Flags().IntVarP(&limit, "limit", "l", 10, "The maximum number of items to return")
	listCmd.Flags().StringVarP(&filter, "filter", "f", "", "Filter items by keyword(s)")
	listCmd.Flags().StringVarP(&properties, "properties", "p", "", "select properties to include in the output: "+helpFmt(&api.ListItem{})+" default:"+helpFmt(&defaultFields))
	listCmd.Flags().StringVarP(&format, "format", "o", "json", "output format. One of: "+helpFmt(&allowedFormats))

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
	rootCmd.AddCommand(Folder())
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

// don't use in big loop, uses reflection
func structToKeys[T any](a *T) []string {
	// Get the reflect.Value of the struct
	v := reflect.ValueOf(*a)

	// Get the reflect.Type of the struct
	t := v.Type()

	var arr []string
	// Iterate over the fields
	for i := 0; i < v.NumField(); i++ {
		// Get the field's Type (contains name, type, tags)
		fieldType := t.Field(i)
		// Get the field's Value
		// fieldValue := v.Field(i)

		jsonName := strings.Split(fieldType.Tag.Get("json"), ",")[0]
		if jsonName != "" {
			arr = append(arr, jsonName)
		} else {
			arr = append(arr, fieldType.Name)
		}

		// fmt.Printf("Field Name: %s, Field Type: %s, Field Value: %v\n",
		//	fieldType.Name,
		//	fieldType.Type,
		//	fieldValue.Interface(), // Use .Interface() to get the actual value
		// )
	}

	return arr
}

func helpFmt[T any](a *T) string {
	val := reflect.Indirect(reflect.ValueOf(a))

	switch val.Kind() {
	case reflect.Slice:
		return strings.Join(val.Interface().([]string), ", ")

	case reflect.Struct:
		// Assuming structToKeys correctly inspects the struct and returns its key names.
		return strings.Join(structToKeys(a), ", ")

	default:
		// Provide a fallback for any other types.
		return fmt.Sprint(val.Interface())
	}
}

// filterStructFields takes a struct and two lists of allowed field names.
// It returns a map containing only the fields that are present in BOTH lists.
func filterStructFields(s any, list1, list2 []string) map[string]any {
	// Use maps for efficient O(1) lookups.
	allowList1 := make(map[string]struct{})
	for _, fieldName := range list1 {
		allowList1[fieldName] = struct{}{}
	}

	allowList2 := make(map[string]struct{})
	for _, fieldName := range list2 {
		allowList2[fieldName] = struct{}{}
	}

	// The map to store the filtered result.
	result := make(map[string]interface{})

	// Use reflection to inspect the struct.
	val := reflect.ValueOf(s)
	// If it's a pointer, get the element it points to.
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	typ := val.Type()
	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		fieldName := field.Name

		// Check if the field name exists in both allow-lists.
		_, inList1 := allowList1[fieldName]
		_, inList2 := allowList2[fieldName]

		if inList1 && inList2 {
			// If it matches, add the field name and value to our result map.
			result[fieldName] = val.Field(i).Interface()
		}
	}

	return result
}

// extractFields takes a struct and a list of desired field names.
// It returns a map containing only the fields from the struct that are present in the list.
func extractFields(s any, fieldsToKeep []string) map[string]any {
	allowList := make(map[string]struct{}, len(fieldsToKeep))
	for _, fieldName := range fieldsToKeep {
		allowList[fieldName] = struct{}{}
	}

	// The map to store the filtered result.
	result := make(map[string]any)

	val := reflect.ValueOf(s)
	val = reflect.Indirect(val)

	// Ensure we are working with a struct.
	if val.Kind() != reflect.Struct {
		return result // Or return an error
	}

	typ := val.Type()
	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		fieldName := field.Name

		// Check if the field name exists in the allow-list.
		if _, allowed := allowList[fieldName]; allowed {
			// If it's allowed, add the field name and value to our result map.
			result[fieldName] = val.Field(i).Interface()
		}
	}

	return result
}

func extractFields1(a any, fields []string) (filteredFields []string) {
	referenceKeys := structToKeys(&a)
	for _, field := range fields {
		for i, _ := range referenceKeys {
			if referenceKeys[i] == field {
				filteredFields = append(filteredFields, field)
			}
		}

	}
	return filteredFields
}
func filterFieldsByReference(referenceFields []string, fields []string) []string {
	referenceSet := make(map[string]struct{})
	for _, field := range referenceFields {
		referenceSet[field] = struct{}{}
	}

	var output []string
	for _, field := range fields {
		if _, found := referenceSet[field]; !found {
			output = append(output, field)
		}
	}

	return output
}
