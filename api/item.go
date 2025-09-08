package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
	_ "net/url"

	"github.com/eissar/nest/api/endpoints"
	// "github.com/eissar/nest/config"
	f "github.com/eissar/nest/format"
	"github.com/spf13/cobra"
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
	Id   string `json:"id"`
	Name string `json:"name"`
	// Size
	// ext
	Tags    []string `json:"tags"`
	Folders []string `json:"folders"`
	// isDeleted
	URL              string `json:"url"`
	Annotation       string `json:"annotation"`
	ModificationTime int64  `json:"modificationTime"`
	// height
	// width
	// lastModified
	// palettes
	Website string `json:"website"`
}

// give a better name
// bulk (addFromUrls) does not include `star` or `folderId`
// pointers represent optional keys and null represents unset
type ItemAddFromUrlOptions struct {
	URL              string            `json:"url"`
	Name             string            `json:"name"`
	Website          string            `json:"website,omitempty"`
	Tags             []string          `json:"tags,omitempty"`
	Star             *int              `json:"star,omitempty"`
	Annotation       string            `json:"annotation,omitempty"`
	ModificationTime int               `json:"modificationTime,omitempty"`
	FolderID         *string           `json:"folderId,omitempty"`
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
	Limit int ` json:"limit"
	flag:"The number of items to be displayed. the default number is 200"`
	Offset  int    `json:"offset,omitempty" flag:"Offset a collection of results from the api. Start with 0."`
	OrderBy string `json:"orderBy,omitempty" flag:"The sorting order.CREATEDATE , FILESIZE , NAME , RESOLUTION , add a minus sign for descending order: -FILESIZE"`
	Keyword string `json:"keyword,omitempty" flag:"Filter by the keyword"`
	Ext     string `json:"ext,omitempty" flag:"Filter by the extension type, e.g.: jpg ,  png"`
	Tags    string `json:"tags,omitempty" flag:"Filter by tags. Use , to divide different tags. E.g.: Design, Poster"`
	Folders string `json:"folders,omitempty" flag:"Filter by Folders.  Use , to divide folder IDs. E.g.: KAY6NTU6UYI5Q,KBJ8Z60O88VMG"`
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

//- [X] /api/item/addFromURL
//- [X] /api/item/addFromURLs
//- [X] /api/item/addFromPath
//- [X] /api/item/addFromPaths
//- [X] /api/item/addBookmark
//- [X] /api/item/info
//- [X] /api/item/thumbnail
//- [X] /api/item/list
//- [X] /api/item/moveToTrash
//- [X] /api/item/refreshPalette
//- [X] /api/item/refreshThumbnail
//- [X] /api/item/update

// endpoint only returns `status`
func ItemAddFromUrl(baseUrl string, item ItemAddFromUrlOptions) error {
	ep := endpoints.ItemAddFromURL
	uri := baseUrl + ep.Path

	body, err := json.Marshal(item)
	if err != nil {
		return fmt.Errorf("addfromurl: error converting request into json body err=%w", err)
	}

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

// endpoint only returns `status`
func ItemAddFromUrls(baseUrl string, items []ItemAddFromUrlOptions, folderId string) error {
	ep := endpoints.ItemAddFromURLs
	uri := baseUrl + ep.Path

	requestBody := struct {
		Items    []ItemAddFromUrlOptions `json:"items"`
		FolderId string                  `json:"folderId,omitempty"`
	}{
		items,
		folderId,
	}
	body, err := json.Marshal(requestBody)
	if err != nil {
		return fmt.Errorf("addfromurls: error converting request into json body err=%w", err)
	}

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
// TODO: endpoint which adds item & returns itemId
func ItemAddFromPath(baseUrl string, item ItemAddFromPathOptions) error {
	ep := endpoints.ItemAddFromPath
	uri := baseUrl + ep.Path

	body, err := json.Marshal(item)
	if err != nil {
		return fmt.Errorf("addfrompath: error converting request into json body err=%w", err)
	}

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

func ItemAddFromPaths(baseUrl string, items []ItemAddFromPathOptions) error {
	ep := endpoints.ItemAddFromPaths
	uri := baseUrl + ep.Path

	requestBody := struct {
		Items []ItemAddFromPathOptions `json:"items"`
		// FolderId string `json:"folderId,omitempty"`
	}{
		items,
		//folderId,
	}

	body, err := json.Marshal(requestBody)
	if err != nil {
		return fmt.Errorf("addfrompaths: error converting request into json body err=%w", err)
	}

	var resp EagleResponse

	err = Request(ep.Method, uri, bytes.NewReader(body), nil, &resp)
	if err != nil {
		return fmt.Errorf("addFromPaths: err=%w", err)
	}

	if resp.Status != "success" {
		return fmt.Errorf("addFromPaths: err=%w", ErrStatusErr)
	}

	return nil
}

// endpoint only returns `status`
func ItemAddBookmark(baseUrl string, item ItemAddBookmarkOptions) error {
	ep := endpoints.ItemAddBookmark
	uri := baseUrl + ep.Path

	// add param checks

	body, err := json.Marshal(item)
	if err != nil {
		return fmt.Errorf("addbookmark: error converting request into json body err=%w", err)
	}

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
func ItemList(baseUrl string, opts ItemListOptions) ([]*ListItem, error) {
	ep := endpoints.ItemList
	uri := baseUrl + ep.Path

	// TODO: validate parameters
	//

	params, err := StructToURLValues(opts)
	if err != nil {
		return nil, fmt.Errorf("list: error converting parameters into url values err=%w", err)
	}

	var resp struct {
		EagleResponse
		Data []*ListItem `json:"data"`
	}

	err = Request(ep.Method, uri, nil, &params, &resp)
	if err != nil {
		return nil, fmt.Errorf("list: err=%w", err)
	}

	if resp.Status != "success" {
		return nil, fmt.Errorf("list: err=%w", ErrStatusErr)
	}

	return resp.Data, nil
}

func ItemMoveToTrash(baseUrl string, ids []string) error {
	ep := endpoints.ItemMoveToTrash
	uri := baseUrl + ep.Path

	// validate itemIds

	respBody := struct {
		ItemIds []string `json:"itemIds"`
	}{
		ids,
	}

	body, err := json.Marshal(respBody)
	if err != nil {
		return fmt.Errorf("itemlist: error converting request into json body err=%w", err)
	}

	resp := EagleResponse{}

	err = Request(ep.Method, uri, bytes.NewReader(body), nil, &resp)
	if err != nil {
		return fmt.Errorf("movetotrash: err=%w", err)
	}

	if resp.Status != "success" {
		return fmt.Errorf("movetotrash: err=%w", ErrStatusErr)
	}

	return nil
}

func ItemRefreshPalette(baseUrl string, id string) error {
	ep := endpoints.ItemRefreshPalette
	uri := baseUrl + ep.Path

	resp := EagleResponse{}

	body, err := json.Marshal(struct {
		Id string `json:"id"`
	}{id})
	if err != nil {
		return fmt.Errorf("itemrefreshpalette: error converting request into json body err=%w", err)
	}

	err = Request(ep.Method, uri, bytes.NewReader(body), nil, &resp)
	if err != nil {
		return fmt.Errorf("itemrefreshpalette: err=%w", err)
	}

	if resp.Status != "success" {
		return fmt.Errorf("itemrefreshpalette: err=%w", ErrStatusErr)
	}

	return nil
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

func ItemRefreshThumbnail(baseUrl string, id string) error {
	ep := endpoints.ItemRefreshPalette
	uri := baseUrl + ep.Path

	resp := EagleResponse{}

	body, err := json.Marshal(struct {
		Id string `json:"id"`
	}{id})
	if err != nil {
		return fmt.Errorf("itemrefreshthumbnail: error converting request into json body err=%w", err)
	}

	err = Request(ep.Method, uri, bytes.NewReader(body), nil, &resp)
	if err != nil {
		return fmt.Errorf("itemrefreshthumbnail: err=%w", err)
	}

	if resp.Status != "success" {
		return fmt.Errorf("itemrefreshthumbnail: err=%w", ErrStatusErr)
	}

	return nil
}

// returns thumbnail path and error
func ItemThumbnail(baseUrl string, id string) (string, error) {
	ep := endpoints.ItemThumbnail
	uri := baseUrl + ep.Path

	// validate query params
	if !IsValidItemID(id) {
		return "", fmt.Errorf("list: error creating request err= itemId parameter malformed or empty.")
	}

	// add query params
	param := url.Values{}
	param.Add("id", id)

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
	if err != nil {
		return respItem, fmt.Errorf("itemupdate: error converting request into json body err=%w", err)
	}

	err = Request(ep.Method, uri, bytes.NewReader(body), nil, &resp)
	if err != nil {
		return respItem, fmt.Errorf("update: err=%w", err)
	}

	if resp.Status != "success" {
		return respItem, fmt.Errorf("update: err=%w", ErrStatusErr)
	}

	return respItem, nil
}

/*
func ItemAddFromUrl(baseUrl string, item ItemAddFromUrlOptions) error {
func ItemAddFromUrls(baseUrl string, items []ItemAddFromUrlOptions, folderId string) error {
func ItemAddFromPath(baseUrl string, item ItemAddFromPathOptions) error {
func ItemAddFromPaths(baseUrl string, items []ItemAddFromPathOptions) error {
func ItemAddBookmark(baseUrl string, item ItemAddBookmarkOptions) error {
func ItemList(baseUrl string, opts ItemListOptions) ([]*ListItem, error) {
func ItemMoveToTrash(baseUrl string, ids []string) error {
func ItemRefreshPalette(baseUrl string, id string) error {
func ItemInfo(baseUrl string, id string) (respItem ApiItem, err error) {
func ItemRefreshThumbnail(baseUrl string, id string) error {
func ItemThumbnail(baseUrl string, itemId string) (string, error) {
func ItemUpdate(baseUrl string, item ItemUpdateOptions) (respItem ApiItem, err error) {
*/

func ItemCmd() *cobra.Command {
	// cfg := config.GetConfig()

	// ! item: `ItemAddFromUrlOptions` | `ItemAddFromPathOptions` | `ItemAddBookmarkOptions` | `ItemUpdateOptions`
	// ! items: `[]ItemAddFromUrlOptions` | `[]ItemAddFromPathOptions`

	// var id string
	// var folderId string
	// var ids []string

	var opts ItemListOptions

	// opts = ItemListOptions{Limit: 10}

	var o f.FormatType

	item := &cobra.Command{
		Use:   "item",
		Short: "Manage items",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(cmd.Flags())
		},
	}

	item.PersistentFlags().VarP(&o, "format", "o", "output format")

	f.BindStructToFlags(item, &opts)

	return item
}
