package api

// - [X] /api/folder/create
// - [X] /api/folder/rename
// - [X] /api/folder/update
// - [X] /api/folder/list
// - [X] /api/folder/listRecent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"

	"github.com/eissar/nest/api/endpoints"
	"github.com/eissar/nest/config"
	f "github.com/eissar/nest/format"
	"github.com/spf13/cobra"
)

// todo rename
type FolderCreateResponse struct {
	ID               string `json:"id"`
	Name             string `json:"name"`
	ModificationTime int    `json:"modificationTime"`
	//Images                     `json:"images"`
	//Folders                    `json:"folders"`
	//ImagesMappings    `json:"imagesMappings"`
	//Tags                       `json:"tags"`
	//Children                   `json:"children"`
	//IsExpand         bool   `json:"isExpand"`
}
type FolderRenameResponse struct {
	ID               string `json:"id"`
	Name             string `json:"name"`
	ModificationTime int    `json:"modificationTime"`
	IsExpand         bool   `json:"isExpand"`
	Size             int    `json:"size"`
	Vstype           string `json:"vstype"`
	IsVisible        bool   `json:"isVisible"`
	HashKey          string `json:"$$hashKey"`
	NewFolderName    string `json:"newFolderName"`
	Editable         bool   `json:"editable"`
	Pinyin           string `json:"pinyin"`
	//Images           []interface{}          `json:"images"`
	//Folders          []interface{}          `json:"folders"`
	//ImagesMappings   map[string]interface{} `json:"imagesMappings"`
	//Tags             []interface{}          `json:"tags"`
	//Children         []interface{}          `json:"children"`
	//Styles           FolderStyles           `json:"styles"`
}
type FolderUpdateResponse struct {
	ID               string `json:"id"`
	Name             string `json:"name"`
	ModificationTime int    `json:"modificationTime"`
	IsExpand         bool   `json:"isExpand"`
	Size             int    `json:"size"`
	Vstype           string `json:"vstype"`
	IsVisible        bool   `json:"isVisible"`
	HashKey          string `json:"$$hashKey"`
	NewFolderName    string `json:"newFolderName"`
	Editable         bool   `json:"editable"`
	Pinyin           string `json:"pinyin"`
	//Images           []interface{}          `json:"images"`
	//Folders          []interface{}          `json:"folders"`
	//ImagesMappings   map[string]interface{} `json:"imagesMappings"`
	//Tags             []interface{}          `json:"tags"`
	//Children         []interface{}          `json:"children"`
	//Styles           FolderStyles           `json:"styles"`
}

// overview information of a folder.
// for some reason, listrecent has more fields...
// pointers are optional fields
type FolderDetailOverview struct {
	ID                   string    `json:"id"`
	Name                 string    `json:"name"`
	Description          string    `json:"description"`
	Children             []any     `json:"children"`
	ModificationTime     int       `json:"modificationTime"`
	Tags                 []string  `json:"tags"`
	Password             *string   `json:"password"`     // optional
	PasswordTips         *string   `json:"passwordTips"` // optional
	Images               *[]string `json:"images"`       // optional
	IsExpand             bool      `json:"isExpand"`
	ImageCount           int       `json:"imageCount"`
	DescendantImageCount int       `json:"descendantImageCount"`
	Pinyin               string    `json:"pinyin"`
	ExtendTags           []string  `json:"extendTags"`
	// ImagesMappings any `json:"imagesMapping"`
	// newFolderName // tf is the point of this key?
}

type FolderStyles struct {
	Depth int  `json:"depth"`
	First bool `json:"first"`
	Last  bool `json:"last"`
}

var folderColors = []string{"red", "orange", "green", "yellow", "aqua", "blue", "purple", "pink"}

// StringSliceContains checks if a string is present in a slice.
func colorsContains(color string) bool {
	for _, c := range folderColors {
		if c == color {
			return true
		}
	}
	return false
}

func FolderCreate(baseUrl string, folderName string) (FolderCreateResponse, error) {
	ep := endpoints.FolderCreate
	uri := baseUrl + ep.Path

	var resp struct {
		EagleResponse
		Data FolderCreateResponse `json:"data"`
	}

	requestBody := struct {
		FolderName string `json:"folderName"`
	}{folderName}

	body, err := json.Marshal(requestBody)
	if err != nil {
		return resp.Data, fmt.Errorf("foldercreate: error converting request into json body err=%w", err)
	}

	err = Request(ep.Method, uri, bytes.NewReader(body), nil, &resp)
	if err != nil {
		return resp.Data, fmt.Errorf("foldercreate: err=%w", err)
	}

	if resp.Status != "success" {
		return resp.Data, fmt.Errorf("foldercreate: err=%w", ErrStatusErr)
	}

	return resp.Data, nil
}
func FolderRename(baseUrl string, folderId string, newName string) /* folder */ error {
	ep := endpoints.FolderRename
	uri := baseUrl + ep.Path

	var resp struct {
		EagleResponse
		Data FolderRenameResponse `json:"data"`
	}

	requestBody := struct {
		FolderId   string `json:"folderId"`
		FolderName string `json:"newName"`
	}{folderId, newName}

	body, err := json.Marshal(requestBody)
	if err != nil {
		return fmt.Errorf("folderrename: error converting request into json body err=%w", err)
	}

	err = Request(ep.Method, uri, bytes.NewReader(body), nil, &resp)
	if err != nil {
		return fmt.Errorf("folderrename: err=%w", err)
	}

	if resp.Status != "success" {
		return fmt.Errorf("folderrename: err=%w", ErrStatusErr)
	}

	return nil
}

// colors
func FolderUpdate(baseUrl string, folderId string, newName string, newDescription string, newColor string) error {
	ep := endpoints.FolderUpdate
	uri := baseUrl + ep.Path

	// validate params
	if newColor != "" {
		if !colorsContains(newColor) {
			return fmt.Errorf("folderupdate: invalid color")
		}
	}
	// validate params

	var resp EagleResponse

	requestBody := struct {
		FolderId       string `json:"folderId"`
		NewName        string `json:"newName,omitempty"`
		NewDescription string `json:"newDescription,omitempty"`
		NewColor       string `json:"newColor,omitempty"`
	}{folderId, newName, newDescription, newColor}

	body, err := json.Marshal(requestBody)
	if err != nil {
		return fmt.Errorf("folderupdate: error converting request into json body err=%w", err)
	}

	err = Request(ep.Method, uri, bytes.NewReader(body), nil, &resp)
	if err != nil {
		return fmt.Errorf("folderupdate: err=%w", err)
	}

	if resp.Status != "success" {
		return fmt.Errorf("folderupdate: err=%w", ErrStatusErr)
	}

	return nil
}

func FolderList(baseUrl string) ([]FolderDetailOverview, error) {
	ep := endpoints.FolderList
	uri := baseUrl + ep.Path

	var resp struct {
		EagleResponse
		Data []FolderDetailOverview `json:"data"`
	}

	err := Request(ep.Method, uri, nil, nil, &resp)
	if err != nil {
		return resp.Data, fmt.Errorf("folderlist: err=%w", err)
	}

	if resp.Status != "success" {
		return resp.Data, fmt.Errorf("folderlist: err=%w", ErrStatusErr)
	}

	return resp.Data, nil
}

func FolderListRecent(baseUrl string) ([]FolderDetailOverview, error) {
	ep := endpoints.FolderListRecent
	uri := baseUrl + ep.Path

	var resp struct {
		EagleResponse
		Data []FolderDetailOverview `json:"data"`
	}

	err := Request(ep.Method, uri, nil, nil, &resp)
	if err != nil {
		return resp.Data, fmt.Errorf("folderlist: err=%w", err)
	}

	if resp.Status != "success" {
		return resp.Data, fmt.Errorf("folderlist: err=%w", ErrStatusErr)
	}

	return resp.Data, nil
}

func addFlags(
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
		Use:   "folder",
		Short: "Manage folders",
	}
	folder.PersistentFlags().VarP(&o, "format", "o", "output format")

	// idFlag := f.CobraPFlagParams{}
	// addIDFlag := func(fs *pflag.FlagSet) {
	// 	fs.StringVarP(&id, idFlag.Name, idFlag.Shorthand, idFlag.Value, idFlag.Usage)
	// }

	folder.AddCommand(
		addFlags(&cobra.Command{
			Use:   "create <name>",
			Short: "Create a new folder",
			Args:  cobra.ExactArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				resp, err := FolderCreate(cfg.BaseURL(), name)
				if err != nil {
					log.Fatalf("FolderCreate: %v", err)
				}
				f.Format(o, resp)
				return nil
			},
		}, nil, &name, nil, nil, nil),
	)
	folder.AddCommand(
		addFlags(&cobra.Command{
			Use:   "rename <id> <new-name>",
			Short: "Rename an existing folder",
			Args:  cobra.ExactArgs(2),
			RunE: func(cmd *cobra.Command, args []string) error {
				if err := FolderRename(cfg.BaseURL(), id, newName); err != nil {
					log.Fatalf("FolderRename: %v", err)
				}
				return nil
			},
		}, &id, nil, &newName, nil, nil),
	)
	folder.AddCommand(addFlags(&cobra.Command{
		Use:   "update <id> <new-name> <new-description> <new-color>",
		Short: "Update folder metadata",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := FolderUpdate(cfg.BaseURL(), args[0], args[1], args[2], args[3]); err != nil {
				log.Fatalf("FolderUpdate: %v", err)
			}
			return nil
		},
	}, &id, nil, &newName, &newDescription, &newColor),
	)
	folder.AddCommand(&cobra.Command{
		Use:   "recent",
		Short: "List recently accessed folders",
		RunE: func(cmd *cobra.Command, args []string) error {
			recent, err := FolderListRecent(cfg.BaseURL())
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
				list, err := FolderList(cfg.BaseURL())
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
