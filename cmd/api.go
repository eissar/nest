// cmd providers from eagle-api
package cmd

import (
	"fmt"
	"log"

	"github.com/eissar/eagle-go"
	"github.com/eissar/nest/config"
	f "github.com/eissar/nest/format"
	"github.com/spf13/cobra"
)

// provides commands
func ApplicationCmd() *cobra.Command {
	cfg := config.GetConfig()

	var o f.FormatType

	app := &cobra.Command{
		Use: "app",
		// Short: "Manage items",
		// Run: func(cmd *cobra.Command, args []string) {
		// 	fmt.Println(cmd.Flags())
		// },
	}
	// return []*cobra.Command{

	func() {
		cmd := &cobra.Command{
			Use:   "info", //
			Short: "Display detailed information about the running Eagle application.",
			Long:  "Retrieves and prints detailed information about the Eagle application currently running. ",
			RunE: func(cmd *cobra.Command, args []string) error {
				v, err := eagle.ApplicationInfo(cfg.BaseURL())

				if err != nil {
					log.Fatalf("Application: %v", err)
				}

				f.Format(o, v)
				return nil
			},
		}
		app.AddCommand(cmd)
	}()
	return app
}

func addFolderFlags(
	cmd *cobra.Command,
	id *string,
	name *string,
	newName *string,
	newDescription *string,
	newColor *string,
) *cobra.Command {
	if id != nil {
		cmd.Flags().StringVar(id, "id", "", "folder id")
	}
	if name != nil {
		cmd.Flags().StringVar(name, "name", "", "folder name")
	}
	if newName != nil {
		cmd.Flags().StringVar(newName, "new-name", "", "updated  folder name")
	}
	if newDescription != nil {
		cmd.Flags().StringVar(newDescription, "description", "", "updated folder description")
	}
	if newColor != nil {
		cmd.Flags().StringVar(newColor, "color", "", "updated folder color")
	}
	return cmd
}

func FolderCmd() *cobra.Command {
	cfg := config.GetConfig()

	var id string
	var name string
	var newName string
	var newDescription string
	var newColor string
	var o f.FormatType

	folder := &cobra.Command{
		Use:   "folder [id]",
		Short: "Manage folders",
		Long:  "use a subcommand or print details for a folder or smart-folder given a single positional ID argument.\nnote: does not find nested smart folders (wip)",

		RunE: func(cmd *cobra.Command, args []string) error {
			// no args
			if len(args) == 0 {
				cmd.Help()
				return nil
			}
			// one positional arg
			if len(args) == 1 {
				targetId := args[0]
				var matchedFolderDetail interface{} // zero value : nil

				cfg := config.GetConfig()
				detail, err := eagle.FolderList(cfg.BaseURL())
				if err != nil {
					return err
				}
				for _, folder := range detail {
					if folder.ID == targetId {
						matchedFolderDetail = folder
						break
					}
				}

				if matchedFolderDetail == nil {
					// we can check if they entered a smart folder
					info, err := eagle.LibraryInfo(cfg.BaseURL())
					if err != nil {
						return err
					}
					for _, folder := range info.SmartFolders {
						if folder.ID == targetId {
							matchedFolderDetail = folder
							break
						}
					}
				}
				if matchedFolderDetail == nil {
					// no matches.
					return fmt.Errorf("Could not find any folder or smart folder matching id (%s)", targetId)
				}

				f.Format(o, matchedFolderDetail)
				return nil
			}

			return nil
		},
	}
	folder.PersistentFlags().VarP(&o, "format", "o", "output format")

	// idFlag := f.CobraPFlagParams{}
	// addIDFlag := func(fs *pflag.FlagSet) {
	// 	fs.StringVarP(&id, idFlag.Name, idFlag.Shorthand, idFlag.Value, idFlag.Usage)
	// }

	folder.AddCommand(
		addFolderFlags(&cobra.Command{
			Use:   "create <name>",
			Short: "Create a new folder",
			Args:  cobra.ExactArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				resp, err := eagle.FolderCreate(cfg.BaseURL(), name)
				if err != nil {
					log.Fatalf("FolderCreate: %v", err)
				}
				f.Format(o, resp)
				return nil
			},
		}, nil, &name, nil, nil, nil),
	)
	folder.AddCommand(
		addFolderFlags(&cobra.Command{
			Use:   "rename <id> <new-name>",
			Short: "Rename an existing folder",
			Args:  cobra.ExactArgs(2),
			RunE: func(cmd *cobra.Command, args []string) error {
				if err := eagle.FolderRename(cfg.BaseURL(), id, newName); err != nil {
					log.Fatalf("FolderRename: %v", err)
				}
				return nil
			},
		}, &id, nil, &newName, nil, nil),
	)
	folder.AddCommand(addFolderFlags(&cobra.Command{
		Use:   "update <id> <new-name> <new-description> <new-color>",
		Short: "Update folder metadata",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := eagle.FolderUpdate(cfg.BaseURL(), id, newName, newDescription, newColor); err != nil {
				log.Fatalf("FolderUpdate: %v", err)
			}
			return nil
		},
	}, &id, nil, &newName, &newDescription, &newColor),
	)
	folder.AddCommand(&cobra.Command{
		Use:   "recent",
		Short: "List recently accessed folders",
		Long:  "List recently accessed folders.\nDoes not enumerate `smart` folders (use nest lib info) for now.",
		RunE: func(cmd *cobra.Command, args []string) error {
			recent, err := eagle.FolderListRecent(cfg.BaseURL())
			if err != nil {
				log.Fatalf("FolderListRecent: %v", err)
			}
			f.Format(o, recent)
			return nil
		}})

	folder.AddCommand(
		&cobra.Command{Use: "list",
			Short: "List all folders",
			RunE: func(cmd *cobra.Command, args []string) error {
				list, err := eagle.FolderList(cfg.BaseURL())
				if err != nil {
					log.Fatalf("FolderList: %v", err)
				}
				f.Format(o, list)
				return nil
			},
		},
	)

	return folder
}

func ItemCmd() *cobra.Command {
	cfg := config.GetConfig()

	// ! items: `[]ItemAddFromUrlOptions` | `[]ItemAddFromPathOptions`

	// var id string
	// var folderId string
	// var ids []string

	// ! item: `ItemAddFromUrlOptions` | `ItemAddFromPathOptions` | `ItemAddBookmarkOptions` | `ItemUpdateOptions`

	// opts = ItemListOptions{Limit: 10}

	var o f.FormatType

	item := &cobra.Command{
		Use: "item",
		// Short: "Manage items",
		// Run: func(cmd *cobra.Command, args []string) {
		// 	fmt.Println(cmd.Flags())
		// },
	}

	item.PersistentFlags().VarP(&o, "format", "o", "output format")

	func() { // [X] use default opts; [X] struct tag metadata
		opts := eagle.ItemListOptions{}.WithDefaults()
		cmd := &cobra.Command{
			Use:   "list",
			Short: "List Items",
			RunE: func(cmd *cobra.Command, args []string) error {
				list, err := eagle.ItemList(cfg.BaseURL(), opts)
				if err != nil {
					log.Fatalf("FolderList: %v", err)
				}
				fmt.Printf("output f: %v\n", o)
				f.Format(o, list)
				return nil
			},
		}

		f.BindStructFlags(cmd, &opts)
		item.AddCommand(cmd)
	}()

	// ItemAddFromUrl
	func() { // [X] use default opts; [X] struct tag metadata
		opts := eagle.ItemAddFromUrlOptions{}.WithDefaults()
		cmd := &cobra.Command{
			Use:   "url [a]",
			Short: "Add item from URL",
			RunE: func(cmd *cobra.Command, args []string) error {
				if err := opts.Validate(); err != nil {
					return err
				}

				err := eagle.ItemAddFromUrl(cfg.BaseURL(), opts)
				if err != nil {
					return fmt.Errorf("failed to add item from URL: %w", err)
				}
				fmt.Println("Successfully added item from URL")
				return nil
			},
		}
		f.BindStructFlags(cmd, &opts)
		item.AddCommand(cmd)
	}()

	// ItemAddFromUrls
	// func() { // [ ] use default opts; [ ] struct tag metadata
	// 	opts := []ItemAddFromUrlOptions{}
	// 	cmd := &cobra.Command{
	// 		Use:   "urls",
	// 		Short: "Add multiple items from URLs",
	// 		RunE: func(cmd *cobra.Command, args []string) error {
	// 			folderId, _ := cmd.Flags().GetString("folder-id")
	// 			err := ItemAddFromUrls(cfg.BaseURL(), opts, folderId)
	// 			if err != nil {
	// 				return fmt.Errorf("failed to add items from URLs: %w", err)
	// 			}
	// 			fmt.Println("Successfully added items from URLs")
	// 			return nil
	// 		},
	// 	}
	// 	f.BindStructFlags(cmd, &opts)
	// 	cmd.Flags().String("folder-id", "", "Folder ID to add items to")
	// 	item.AddCommand(cmd)
	// }()

	// ItemAddFromPath
	func() { // [X] use default opts; [X] struct tag metadata
		opts := eagle.ItemAddFromPathOptions{}
		cmd := &cobra.Command{
			Use:   "path",
			Short: "Add item from local path",
			RunE: func(cmd *cobra.Command, args []string) error {
				if err := opts.Validate(); err != nil {
					return err
				}

				err := eagle.ItemAddFromPath(cfg.BaseURL(), opts)
				if err != nil {
					return fmt.Errorf("failed to add item from path: %w", err)
				}
				fmt.Println("Successfully added item from path")
				return nil
			},
		}
		f.BindStructFlags(cmd, &opts)
		item.AddCommand(cmd)
	}()

	// ItemAddFromPaths
	// func() { // [ ] use default opts; [ ] struct tag metadata
	// 	opts := []ItemAddFromPathOptions{}
	// 	cmd := &cobra.Command{
	// 		Use:   "paths",
	// 		Short: "Add multiple items from local paths",
	// 		RunE: func(cmd *cobra.Command, args []string) error {
	// 			err := ItemAddFromPaths(cfg.BaseURL(), opts)
	// 			if err != nil {
	// 				return fmt.Errorf("failed to add items from paths: %w", err)
	// 			}
	// 			fmt.Println("Successfully added items from paths")
	// 			return nil
	// 		},
	// 	}
	// 	f.BindStructFlags(cmd, &opts)
	// 	item.AddCommand(cmd)
	// }()

	// ItemAddBookmark
	func() { // [X] use default opts; [X] struct tag metadata
		opts := eagle.ItemAddBookmarkOptions{}.WithDefaults()
		cmd := &cobra.Command{
			Use:   "bookmark",
			Short: "Add bookmark item",
			RunE: func(cmd *cobra.Command, args []string) error {
				if err := opts.Validate(); err != nil {
					return err
				}

				err := eagle.ItemAddBookmark(cfg.BaseURL(), opts)

				if err != nil {
					return fmt.Errorf("failed to add bookmark: %w", err)
				}
				fmt.Println("Successfully added bookmark")
				return nil
			},
		}
		f.BindStructFlags(cmd, &opts)
		item.AddCommand(cmd)
	}()

	// ItemMoveToTrash
	func() {
		ids := []string{}
		cmd := &cobra.Command{
			Use:   "delete",
			Short: "Move item to the trash",
			RunE: func(cmd *cobra.Command, args []string) error {
				if len(ids) == 0 && len(args) > 0 {
					ids = args
				}
				err := eagle.ItemMoveToTrash(cfg.BaseURL(), ids)
				if err != nil {
					return fmt.Errorf("failed to move items to trash: %w", err)
				}
				fmt.Printf("Successfully moved %d items to trash\n", len(ids))
				return nil
			},
		}
		cmd.Flags().StringSliceVar(&ids, "ids", []string{}, "Item IDs to move to trash")
		item.AddCommand(cmd)
	}()

	// ItemRefreshPalette
	func() {
		var id string
		cmd := &cobra.Command{
			Use:   "refresh-palette",
			Short: "Refresh item color palette",
			RunE: func(cmd *cobra.Command, args []string) error {
				if id == "" && len(args) > 0 {
					id = args[0]
				}
				err := eagle.ItemRefreshPalette(cfg.BaseURL(), id)
				if err != nil {
					return fmt.Errorf("failed to refresh palette: %w", err)
				}
				fmt.Println("Successfully refreshed color palette")
				return nil
			},
		}
		cmd.Flags().StringVar(&id, "id", "", "Item ID")
		item.AddCommand(cmd)
	}()

	// ItemInfo
	func() {
		var id string
		cmd := &cobra.Command{
			Use:   "info",
			Short: "Get item info",
			RunE: func(cmd *cobra.Command, args []string) error {
				if id == "" && len(args) > 0 {
					id = args[0]
				}
				resp, err := eagle.ItemInfo(cfg.BaseURL(), id)
				if err != nil {
					return fmt.Errorf("failed to get item info: %w", err)
				}
				f.Format(o, resp)
				return nil
			},
		}
		cmd.Flags().StringVar(&id, "id", "", "Item ID")
		item.AddCommand(cmd)
	}()

	// ItemRefreshThumbnail
	func() {
		var id string
		cmd := &cobra.Command{
			Use:   "refresh-thumbnail",
			Short: "Refresh item thumbnail",
			RunE: func(cmd *cobra.Command, args []string) error {
				if id == "" && len(args) > 0 {
					id = args[0]
				}
				err := eagle.ItemRefreshThumbnail(cfg.BaseURL(), id)
				if err != nil {
					return fmt.Errorf("failed to refresh thumbnail: %w", err)
				}
				fmt.Println("Successfully refreshed thumbnail")
				return nil
			},
		}
		cmd.Flags().StringVar(&id, "id", "", "Item ID")
		item.AddCommand(cmd)
	}()

	// ItemThumbnail
	func() {
		var itemId string
		cmd := &cobra.Command{
			Use:   "thumbnail",
			Short: "Get item thumbnail",
			RunE: func(cmd *cobra.Command, args []string) error {
				if itemId == "" && len(args) > 0 {
					itemId = args[0]
				}
				thumbnail, err := eagle.ItemThumbnail(cfg.BaseURL(), itemId)
				if err != nil {
					return fmt.Errorf("failed to get thumbnail: %w", err)
				}
				fmt.Println(thumbnail)
				return nil
			},
		}
		cmd.Flags().StringVar(&itemId, "id", "", "Item ID")
		item.AddCommand(cmd)
	}()

	// ItemUpdate
	func() {
		opts := eagle.ItemUpdateOptions{}
		cmd := &cobra.Command{
			Use:   "update",
			Short: "Update item",
			RunE: func(cmd *cobra.Command, args []string) error {
				if err := opts.Validate(); err != nil {
					return err
				}

				resp, err := eagle.ItemUpdate(cfg.BaseURL(), opts)
				if err != nil {
					return fmt.Errorf("failed to update item: %w", err)
				}
				f.Format(o, resp)
				return nil
			},
		}
		f.BindStructFlags(cmd, &opts)
		item.AddCommand(cmd)
	}()

	return item
}

func LibraryCmd() *cobra.Command {
	cfg := config.GetConfig()

	var o f.FormatType

	library := &cobra.Command{
		Use:   "lib",
		Short: "Manage Libraries",
	}
	library.PersistentFlags().VarP(&o, "format", "o", "output format")

	func() { // LibraryInfo
		cmd := &cobra.Command{
			Use:   "info",
			Short: "Display current library details",
			Args:  cobra.NoArgs,
			RunE: func(cmd *cobra.Command, args []string) error {
				data, err := eagle.LibraryInfo(cfg.BaseURL())
				if err != nil {
					return err
				}
				f.Format(o, data)
				return nil
			},
		}

		library.AddCommand(cmd)
	}()

	func() { // LibraryHistory
		cmd := &cobra.Command{
			Use:   "history",
			Short: "List libraries in the recent list in the menu bar > libraries",
			Args:  cobra.NoArgs,
			RunE: func(cmd *cobra.Command, args []string) error {
				data, err := eagle.LibraryHistory(cfg.BaseURL())
				if err != nil {
					return err
				}
				f.Format(o, data)
				return nil
			},
		}
		library.AddCommand(cmd)
	}()

	func() { // LibrarySwitch
		var libraryPath string

		cmd := &cobra.Command{
			Use:   "switch [library-path]",
			Short: "Change active library to the given path",
			Args:  cobra.MaximumNArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				path := libraryPath
				if len(args) > 0 {
					path = args[0]
				}
				if path == "" {
					return fmt.Errorf("library path is required")
				}
				return eagle.LibrarySwitch(cfg.BaseURL(), path)
			},
		}
		cmd.Flags().StringVarP(&libraryPath, "librarypath", "L", "", "path to library")
		library.AddCommand(cmd)
	}()

	// TODO: this endpoint always returns `library does not exist?`
	// library.AddCommand(
	// 	&cobra.Command{
	// 		Use:   "icon",
	// 		Short: "Return the URL of the current library’s icon",
	// 		// Args:  cobra.NoArgs,
	// 		RunE: func(cmd *cobra.Command, args []string) error {
	// 			url, err := LibraryIcon(cfg.BaseURL())
	// 			if err != nil {
	// 				return err
	// 			}
	// 			cmd.Println(url)
	// 			return nil
	// 		},
	// 	},
	// )

	return library
}
