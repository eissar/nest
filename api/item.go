package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	_ "net/url"
	"strconv"

	"github.com/eissar/nest/api/endpoints"
)

// #region types

// site, annotation, tags, folderid ?
type BaseItem struct {
}
type Palette struct {
	Color   []int   `json:"color"`
	Ratio   float64 `json:"ratio"`
	HashKey string  `json:"$$hashKey"`
}
type Item struct {
	URL              string   `json:"url"`
	Name             string   `json:"name"`
	Website          string   `json:"website"`
	Tags             []string `json:"tags"`
	Star             int      `json:"star"`
	Annotation       string   `json:"annotation"`
	ModificationTime int64    `json:"modificationTime"`
	FolderID         string   `json:"folderId"`
	//Headers          map[string]string `json:"headers,omitempty"`
}

// Fields returned by Eagle API endpoints (item/list, item/info, item/update).
//
// NOTE:
//   - item/info includes the optional key `noThumbnail`.
//   - item/update includes `noThumbnail` and the optional key `star`.
//   - For optional keys null represents unset
type ApiItem struct {
	ID               string    `json:"id"`
	Name             string    `json:"name"`
	Size             int       `json:"size"`
	Ext              string    `json:"ext"`
	Tags             []string  `json:"tags"`
	Folders          []string  `json:"folders"`
	IsDeleted        bool      `json:"isDeleted"`
	URL              string    `json:"url"`
	Annotation       string    `json:"annotation"`
	ModificationTime int64     `json:"modificationTime"`
	Width            int       `json:"width"`
	Height           int       `json:"height"`
	NoThumbnail      *bool     `json:"noThumbnail,omitempty"`
	LastModified     int64     `json:"lastModified"`
	Palettes         []Palette `json:"palettes"`
	Star             *int      `json:"star,omitempty"`
}

type ListItem struct {
	Id               string   `json:"id"`
	Name             string   `json:"name"`
	URL              string   `json:"url"`
	Website          string   `json:"website"`
	Tags             []string `json:"tags"`
	Annotation       string   `json:"annotation"`
	ModificationTime int64    `json:"modificationTime"`
	//FolderID         string   `json:"folderId"`
}

// give a better name
type ItemAddFromUrlOptions struct {
	URL              string            `json:"url"`
	Name             string            `json:"name"`
	Website          string            `json:"website,omitempty"`
	Tags             []string          `json:"tags,omitempty"`
	Star             int               `json:"star,omitempty"`
	Annotation       string            `json:"annotation,omitempty"`
	ModificationTime int               `json:"modificationTime,omitempty"`
	FolderID         string            `json:"folderId,omitempty"`
	Headers          map[string]string `json:"headers,omitempty"`
}

// give better name
type ItemAddBookmarkOptions struct {
	URL              string   `json:"url"`
	Name             string   `json:"name"`
	Base64           string   `json:"base64,omitempty"`
	Tags             []string `json:"tags,omitempty"`
	ModificationTime string   `json:"modificationTime,omitempty"`
	FolderID         string   `json:"folderId,omitempty"`
}

// pointers represent optional keys and null represents unset
type ItemUpdateOptions struct {
	ID         string    `json:"id"`
	Tags       *[]string `json:"tags"`
	Annotation *string   `json:"annotation"`
	URL        *string   `json:"url"`
	Star       *int      `json:"star"`
}

// no folder Id
type BulkItem struct {
	Item
	//FolderId string `json:"omitempty`
}
type ItemListOptions struct {
	Limit   int    `json:"limit"`             // The number of items to be displayed. the default number is 200
	Offset  int    `json:"offset,omitempty"`  // Offset a collection of results from the api. Start with 0.
	OrderBy string `json:"orderBy,omitempty"` // The sorting order.CREATEDATE , FILESIZE , NAME , RESOLUTION , add a minus sign for descending order: -FILESIZE
	Keyword string `json:"keyword,omitempty"` // Filter by the keyword
	Ext     string `json:"ext,omitempty"`     // Filter by the extension type, e.g.: jpg ,  png
	Tags    string `json:"tags,omitempty"`    // Filter by tags. Use , to divide different tags. E.g.: Design, Poster
	Folders string `json:"folders,omitempty"` // Filter by Folders.  Use , to divide folder IDs. E.g.: KAY6NTU6UYI5Q,KBJ8Z60O88VMG
}
type ItemAddFromPathOptions struct {
	Path       string   `json:"path"`                 // Required, the path of the local file.
	Name       string   `json:"name,omitempty"`       // Required, the name of the image to be added. (not really req)
	Website    string   `json:"website,omitempty"`    // The Address of the source of the image.
	Annotation string   `json:"annotation,omitempty"` // The annotation for the image.
	Tags       []string `json:"tags,omitempty"`       // Tags for the image.
	FolderId   string   `json:"folderId,omitempty"`   // If this parameter is defined, the image will be added to the corresponding folder.
}
type ThumbnailData struct {
	Status        string `json:"status"`
	ThumbnailPath string `json:"data"`
}

// #endregion types

// start api endpoints

// use cleaner abstraction `api.Request` like in ItemInfo

//- [X] /api/item/addFromURL
//- [ ] /api/item/addFromURLs
//- [X] /api/item/addFromPath
//- [ ] /api/item/addFromPaths
//- [X] /api/item/addBookmark
//- [X] /api/item/info
//- [X] /api/item/thumbnail
//- [X] /api/item/list
//- [ ] /api/item/moveToTrash
//- [ ] /api/item/refreshPalette
//- [ ] /api/item/refreshThumbnail
//- [X] /api/item/update

// endpoint only returns `status`
func ItemAddFromUrl(baseUrl string, item ItemAddFromUrlOptions) error {
	ep := endpoints.ItemAddFromURL
	uri := baseUrl + ep.Path

	body, err := json.Marshal(item)

	var resp EagleResponse

	err = Request(ep.Method, uri, bytes.NewReader(body), nil, &resp)
	if err != nil {
		return fmt.Errorf("addFromUrl: err=%w", err)
	}

	if resp.Status != "success" {
		return fmt.Errorf("addFromUrl: err=%w", ErrStatusErr)
	}

	return nil
}

// returns status only.
// use NewItemFromPath to build args
func ItemAddFromPath(baseURL string, item ItemAddFromPathOptions) error {
	ep := endpoints.ItemAddFromPath
	uri := baseURL + ep.Path

	body, err := json.Marshal(item)

	var resp EagleResponse

	err = Request(ep.Method, uri, bytes.NewReader(body), nil, &resp)
	if err != nil {
		return fmt.Errorf("addFromPath: err=%w", err)
	}

	if resp.Status != "success" {
		return fmt.Errorf("addFromPath: err=%w", ErrStatusErr)
	}

	return nil
}

// endpoint only returns `status`
func ItemAddBookmark(baseUrl string, item ItemAddBookmarkOptions) error {
	ep := endpoints.ItemAddBookmark
	uri := baseUrl + ep.Path

	// add param checks

	body, err := json.Marshal(item)

	var resp EagleResponse

	err = Request(ep.Method, uri, bytes.NewReader(body), nil, &resp)
	if err != nil {
		return fmt.Errorf("AddBookmark: err=%w", err)
	}

	if resp.Status != "success" {
		return fmt.Errorf("AddBookmark: err=%w", ErrStatusErr)
	}

	return nil
}

// creates an *http.Request and sends to InvokeEagleAPIV1
func ItemList(baseUrl string, opts ItemListOptions) ([]ListItem, error) {
	ep := endpoints.ItemList
	uri := baseUrl + ep.Path

	// TODO: validate parameters
	//
	// param
	param := url.Values{}
	if opts.Limit > 0 {
		param.Add("limit", strconv.Itoa(opts.Limit))
	}

	// end param

	var resp struct {
		EagleResponse
		Data []ListItem `json:"data"`
	}

	err := Request(ep.Method, uri, nil, &param, &resp)
	if err != nil {
		return nil, fmt.Errorf("list: err=%w", err)
	}

	if resp.Status != "success" {
		return nil, fmt.Errorf("list: err=%w", ErrStatusErr)
	}

	return resp.Data, nil
}

func ItemInfo(baseUrl string, id string) (respItem ApiItem, err error) {
	//#region Validate
	if !IsValidItemID(id) {
		return respItem, fmt.Errorf("list: error creating request err= itemId parameter malformed or empty.")
	}
	//#endregion Validate

	ep := endpoints.ItemInfo
	uri := baseUrl + ep.Path

	param := url.Values{}
	param.Add("id", id)

	var resp struct {
		EagleResponse
		Data ApiItem `json:"data"`
	}
	err = Request(ep.Method, uri, nil, &param, &resp)
	if err != nil {
		return respItem, fmt.Errorf("ItemInfo: err=%w", err)
	}

	return resp.Data, nil
}

func ItemUpdate(baseUrl string, item ItemUpdateOptions) (respItem ApiItem, err error) {
	ep := endpoints.ItemUpdate
	uri := baseUrl + ep.Path

	// validate id
	if !IsValidItemID(item.ID) {
		return respItem, fmt.Errorf("list: error creating request err= itemId parameter malformed or empty.")
	}
	//end validations

	var resp struct {
		EagleResponse
		Data ApiItem `json:"data"`
	}

	body, err := json.Marshal(item)

	err = Request(ep.Method, uri, bytes.NewReader(body), nil, &resp)
	if err != nil {
		return respItem, fmt.Errorf("update: err=%w", err)
	}

	if resp.Status != "success" {
		return respItem, fmt.Errorf("update: err=%w", ErrStatusErr)
	}

	return respItem, nil
}

// returns thumbnail path and error
func ItemThumbnail(baseUrl string, itemId string) (string, error) {
	ep := endpoints.ItemThumbnail
	uri := baseUrl + ep.Path

	// validate query params
	if !IsValidItemID(itemId) {
		return "", fmt.Errorf("list: error creating request err= itemId parameter malformed or empty.")
	}

	// add query params
	param := url.Values{}
	param.Add("id", itemId)

	// TODO: replace with struct
	var resp ThumbnailData

	err := Request(ep.Method, uri, nil, &param, &resp)
	if err != nil {
		return "", fmt.Errorf("thumbnail: err=%w", err)
	}

	if resp.Status != "success" {
		return "", fmt.Errorf("update: err=%w", ErrStatusErr)
	}

	if escapedPath, err := url.PathUnescape(resp.ThumbnailPath); err != nil {
		return resp.ThumbnailPath, fmt.Errorf("could not url decode path response from eagle server err=%w", err)
	} else {
		return escapedPath, nil
	}
}
