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

	"github.com/eissar/nest/api/endpoints"
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
	Children             []string  `json:"children"`
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
