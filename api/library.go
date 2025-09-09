package api

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/eissar/nest/api/endpoints"
	"github.com/eissar/nest/config"
	f "github.com/eissar/nest/format"
	"github.com/spf13/cobra"
)

// #region types

type Folder struct {
	ID               string   `json:"id"`
	Name             string   `json:"name"`
	Description      string   `json:"description"`
	Children         []Folder `json:"children"`
	ModificationTime int64    `json:"modificationTime,omitempty"`
	Tags             []string `json:"tags,omitempty"`
	IconColor        string   `json:"iconColor,omitempty"`
	Password         string   `json:"password,omitempty"`
	PasswordTips     string   `json:"passwordTips,omitempty"`
	CoverID          string   `json:"coverId,omitempty"`
	OrderBy          string   `json:"orderBy,omitempty"`
	SortIncrease     bool     `json:"sortIncrease,omitempty"`
	Icon             string   `json:"icon,omitempty"`
}
type SmartFolder struct {
	ID               string      `json:"id"`
	Icon             string      `json:"icon"`
	Name             string      `json:"name"`
	Description      string      `json:"description"`
	ModificationTime int64       `json:"modificationTime"`
	Conditions       []Condition `json:"conditions"`
	OrderBy          string      `json:"orderBy,omitempty"`
	SortIncrease     bool        `json:"sortIncrease,omitempty"`
}
type Library struct {
	Path string `json:"path"`
	Name string `json:"name"`
}
type LibraryData struct {
	Folders            []Folder      `json:"folders"`
	SmartFolders       []SmartFolder `json:"smartFolders"`
	QuickAccess        []QuickAccess `json:"quickAccess"`
	TagsGroups         []TagsGroup   `json:"tagsGroups"`
	ModificationTime   int64         `json:"modificationTime"`
	ApplicationVersion string        `json:"applicationVersion"`
	Library            Library       `json:"library"`
}
type Condition struct {
	HashKey string `json:"$$hashKey,omitempty"`
	Match   string `json:"match"`
	Rules   []Rule `json:"rules"`
}
type Rule struct {
	HashKey  string `json:"$$hashKey,omitempty"`
	Method   string `json:"method"`
	Property string `json:"property"`
	Value    any    `json:"value"` // Can be []int or string or []string
}
type QuickAccess struct {
	Type string `json:"type"`
	ID   string `json:"id"`
}
type TagsGroup struct {
	ID    string   `json:"id"`
	Name  string   `json:"name"`
	Tags  []string `json:"tags"`
	Color string   `json:"color,omitempty"`
}
type LibraryInfoResponse struct {
	Data   LibraryData `json:"data"`
	Status string      `json:"status"`
}

// #endregion types

// start endpoints

//- [X] /api/library/info
//- [X] /api/library/history
//- [X] /api/library/switch
//- [-] /api/library/icon

func LibraryInfo(baseURL string) (*LibraryData, error) {
	ep := endpoints.LibraryInfo
	uri := baseURL + ep.Path

	var resp struct {
		EagleResponse             // `json:"response"`
		Data          LibraryData `json:"data"`
	}

	err := Request(ep.Method, uri, nil, nil, &resp)
	if err != nil {
		return &resp.Data, fmt.Errorf("LibraryInfo: err=%w", err)
	}

	if resp.Status != "success" {
		return &resp.Data, fmt.Errorf("LibraryInfo: err=%w", ErrStatusErr)
	}

	return &resp.Data, nil
}

// returns []string paths to libraries
// /api/library/history
func LibraryHistory(baseURL string) ([]string, error) {
	ep := endpoints.LibraryHistory
	uri := baseURL + ep.Path

	var resp struct {
		EagleResponse
		Data []string `json:"data"`
	}

	err := Request(ep.Method, uri, nil, nil, &resp)
	if err != nil {
		return []string{}, fmt.Errorf("recent: err=%w", err)
	}

	if resp.Status != "success" {
		return []string{}, fmt.Errorf("recent: err=%w", ErrStatusErr)
	}

	return resp.Data, nil
}

// cleans libraryPath and tries to switch.
// /api/library/switch
// endpoint only returns `status`
func LibrarySwitch(baseURL string, libraryPath string) error {
	ep := endpoints.LibrarySwitch
	uri := baseURL + ep.Path

	// validate params

	if _, err := os.Stat(libraryPath); err != nil {
		return fmt.Errorf("switch: err=%w", err)
	}

	libraryPath = filepath.Clean(libraryPath)
	libraryPath = filepath.ToSlash(libraryPath)
	libraryPath = strings.TrimSuffix(libraryPath, "/") // issue ...

	// end validate params

	var resp EagleResponse

	body := fmt.Appendf(nil, `{"libraryPath": "%s"}`, libraryPath) // bytes

	err := Request(ep.Method, uri, bytes.NewReader(body), nil, &resp)
	if err != nil {
		return fmt.Errorf("switch: err=%w", err)
	}

	if resp.Status != "success" {
		return fmt.Errorf("switch: err=%w", ErrStatusErr)
	}

	return nil
}

// returns string iconpath (broken)
func LibraryIcon(baseURL string) (string, error) {
	var currentLibraryPath string

	ep := endpoints.LibraryIcon

	uri := baseURL + ep.Path

	req, err := http.NewRequest(ep.Method, uri, nil) // method, url, body
	if err != nil {
		return currentLibraryPath, fmt.Errorf("getcurrlibrarypath: error creating request err=%w", err)
	}

	// FIX
	var a *EagleMessage
	err = invokeEagleAPI(req, &a)
	if err != nil {
		return currentLibraryPath, fmt.Errorf("getcurrlibrarypath: error invoking request err=%w", err)
	}

	if v, ok := a.Data.(string); ok {
		currentLibraryPath, err = url.PathUnescape(v)
		if err != nil {
			return currentLibraryPath, fmt.Errorf("getcurrlibrarypath: error parsing path err=%w", err)
		}
	}

	return currentLibraryPath, nil
}

func addLibFlags(cmd *cobra.Command, libraryPath *string) *cobra.Command {
	if libraryPath != nil {
		cmd.Flags().StringVarP(libraryPath, "librarypath", "L", "", "path to library")
	}
	return cmd
}
func LibraryCmd() *cobra.Command {
	cfg := config.GetConfig()

	var o f.FormatType
	var libraryPath string

	library := &cobra.Command{
		Use:   "lib",
		Short: "Manage Libraries",
	}
	library.PersistentFlags().VarP(&o, "format", "o", "output format")

	library.AddCommand(
		&cobra.Command{
			Use:   "info",
			Short: "Display current library details",
			Args:  cobra.NoArgs,
			RunE: func(cmd *cobra.Command, args []string) error {
				data, err := LibraryInfo(cfg.BaseURL())
				if err != nil {
					return err
				}
				f.Format(o, data)
				return nil
			},
		},
	)

	library.AddCommand(
		&cobra.Command{
			Use:   "history",
			Short: "List libraries in the recent list in the menu bar > libraries",
			Args:  cobra.NoArgs,
			RunE: func(cmd *cobra.Command, args []string) error {
				data, err := LibraryHistory(cfg.BaseURL())
				if err != nil {
					return err
				}
				f.Format(o, data)
				return nil
			},
		},
	)

	library.AddCommand(
		addLibFlags(&cobra.Command{
			Use:   "switch [library-path]",
			Short: "Change active library to the given path",
			Args:  cobra.ExactArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				return LibrarySwitch(cfg.BaseURL(), libraryPath)
			},
		}, &libraryPath),
	)

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
