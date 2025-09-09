package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	_ "net/url"
	"os"

	"github.com/eissar/nest/api/endpoints"
	"github.com/eissar/nest/config"

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
// defaults to use url as item name in eagle
// bulk (addFromUrls) does not include `star` or `folderId`
// pointers represent optional keys and null represents unset
type ItemAddFromUrlOptions struct {
	URL              string            `json:"url" flagname:"u" flag:"url to item to add"`
	Name             string            `json:"name" flag:"name to use for item"`
	Website          string            `json:"website,omitempty" flag:"associated website of item"`
	Tags             []string          `json:"tags,omitempty" flag:"tags to apply to item"`
	Star             *int              `json:"star,omitempty" flag:"star rating of the item"`
	Annotation       string            `json:"annotation,omitempty" flag:"annotation text for the item"`
	ModificationTime int               `json:"modificationTime,omitempty" flag:"modification time in epoch milliseconds"`
	FolderID         *string           `json:"folderId,omitempty" flag:"folder id to place the item in"`
	Headers          map[string]string `json:"headers,omitempty" flag:"http headers to be sent with requests"`
}

func (o ItemAddFromUrlOptions) WithDefaults() (ItemAddFromUrlOptions, error) {
	if o.URL == "" {
		return o, fmt.Errorf("ItemAddFromUrlOptions: url is required")
	}
	if o.Name == "" {
		o.Name = o.URL
	}
	return o, nil
}

// give better name
type ItemAddBookmarkOptions struct {
	URL string `json:"url" flag:"URL of the bookmark"`

	Name             string   `json:"name" flag:"Display name for the bookmark"`
	Base64           string   `json:"base64,omitempty" flag:"Optional base64-encoded data"`
	Tags             []string `json:"tags,omitempty" flag:"Optional list of tag names"`
	ModificationTime string   `json:"modificationTime,omitempty" flag:"DESC TODO"`
	FolderID         string   `json:"folderId,omitempty" flag:"Optional ID of target folder to place the bookmark"`
}

func (o ItemAddBookmarkOptions) WithDefaults() (ItemAddBookmarkOptions, error) {
	if o.URL == "" {
		return o, fmt.Errorf("ItemAddBookmarkOptions: url is required")
	}
	if o.Name == "" {
		o.Name = o.URL
	}
	return o, nil
}

// pointers represent optional keys and null represents unset
type ItemUpdateOptions struct {
	ID         string    `json:"id" flag:"unique identifier"`
	Tags       *[]string `json:"tags" flag:"list of tags associated with the item"`
	Annotation *string   `json:"annotation" flag:"user-provided annotation or note"`
	URL        *string   `json:"url" flag:"web URL associated with the item"`
	Star       *int      `json:"star" flag:"star rating from 1-5, nil for no rating"`
}

func (o ItemUpdateOptions) WithDefaults() (ItemUpdateOptions, error) {
	if o.ID == "" {
		return o, fmt.Errorf("ItemUpdateOptions: id is required")
	}
	if o.Tags == nil && o.Annotation == nil && o.URL == nil && o.Star == nil {
		return o, fmt.Errorf("ItemUpdateOptions: no updates specified - at least one field must be set")
	}
	return o, nil
}

// no folder Id
type BulkItem struct {
	Item
	//FolderId string `json:"omitempty`
}

type ItemListOptions struct {
	Limit   int    `json:"limit" flag:"The number of items to be displayed. the default number is 200"`
	Offset  int    `json:"offset,omitempty" flag:"Offset a collection of results from the api. Start with 0."`
	OrderBy string `json:"orderBy,omitempty" flag:"The sorting order. CREATEDATE , FILESIZE , NAME , RESOLUTION , add a minus sign for descending order: -FILESIZE"`
	Keyword string `json:"keyword,omitempty" flag:"Filter by the keyword"`
	Ext     string `json:"ext,omitempty" flag:"Filter by the extension type, e.g.: jpg ,  png"`
	Tags    string `json:"tags,omitempty" flag:"Filter by tags. Use , to divide different tags. E.g.: Design, Poster"`
	Folders string `json:"folders,omitempty" flag:"Filter by Folders.  Use , to divide folder IDs. E.g.: KAY6NTU6UYI5Q,KBJ8Z60O88VMG"`
}

func (o ItemListOptions) WithDefaults() ItemListOptions {
	if o.Limit == 0 {
		o.Limit = 200
	}
	return o
}

type ItemAddFromPathOptions struct {
	Path       string   `json:"path" flag:"Required, the path of the local file."`
	Name       string   `json:"name,omitempty" flag:"Required, the name of the image to be added. (not really req)"`
	Website    string   `json:"website,omitempty" flag:"The Address of the source of the image."`
	Annotation string   `json:"annotation,omitempty" flag:"The annotation for the image."`
	Tags       []string `json:"tags,omitempty" flag:"Tags for the image."`
	FolderId   string   `json:"folderId,omitempty" flag:"If this parameter is defined, the image will be added to the corresponding folder."`
}

// WithDefaults returns a validated copy of the options with all defaults applied.
// It checks that Path is set and valid.
func (o ItemAddFromPathOptions) WithDefaults() (ItemAddFromPathOptions, error) {
	if o.Path == "" {
		return o, fmt.Errorf("ItemAddFromPathOptions: path is required")
	}
	if _, err := os.Stat(o.Path); err != nil {
		return o, fmt.Errorf("ItemAddFromPathOptions: invalid path: %w", err)
	}
	return o, nil
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
[ ] func ItemAddFromUrl(baseUrl string, item ItemAddFromUrlOptions) error {
[ ] func ItemAddFromUrls(baseUrl string, items []ItemAddFromUrlOptions, folderId string) error {
[ ] func ItemAddFromPath(baseUrl string, item ItemAddFromPathOptions) error {
[ ] func ItemAddFromPaths(baseUrl string, items []ItemAddFromPathOptions) error {
[ ] func ItemAddBookmark(baseUrl string, item ItemAddBookmarkOptions) error {
[ ] func ItemList(baseUrl string, opts ItemListOptions) ([]*ListItem, error) {
[ ] func ItemMoveToTrash(baseUrl string, ids []string) error {
[ ] func ItemRefreshPalette(baseUrl string, id string) error {
[ ] func ItemInfo(baseUrl string, id string) (respItem ApiItem, err error) {
[ ] func ItemRefreshThumbnail(baseUrl string, id string) error {
[ ] func ItemThumbnail(baseUrl string, itemId string) (string, error) {
[ ] func ItemUpdate(baseUrl string, item ItemUpdateOptions) (respItem ApiItem, err error) {
[ ] */

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
		opts := ItemListOptions{}.WithDefaults()
		cmd := &cobra.Command{
			Use:   "list",
			Short: "List Items",
			RunE: func(cmd *cobra.Command, args []string) error {
				list, err := ItemList(cfg.BaseURL(), opts)
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
		opts, defaultErr := ItemAddFromUrlOptions{}.WithDefaults()
		cmd := &cobra.Command{
			Use:   "url [a]",
			Short: "Add item from URL",
			RunE: func(cmd *cobra.Command, args []string) error {
				if defaultErr != nil {
					return defaultErr // should be good enough
				}

				err := ItemAddFromUrl(cfg.BaseURL(), opts)
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
		opts, defaultErr := ItemAddFromPathOptions{}.WithDefaults()
		cmd := &cobra.Command{
			Use:   "path",
			Short: "Add item from local path",
			RunE: func(cmd *cobra.Command, args []string) error {
				if defaultErr != nil {
					return defaultErr // should be good enough
				}
				err := ItemAddFromPath(cfg.BaseURL(), opts)
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
		opts, defaultErr := ItemAddBookmarkOptions{}.WithDefaults()
		cmd := &cobra.Command{
			Use:   "bookmark",
			Short: "Add bookmark item",
			RunE: func(cmd *cobra.Command, args []string) error {
				if defaultErr != nil {
					return defaultErr // should be good enough
				}
				err := ItemAddBookmark(cfg.BaseURL(), opts)
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
				err := ItemMoveToTrash(cfg.BaseURL(), ids)
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
				err := ItemRefreshPalette(cfg.BaseURL(), id)
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
				resp, err := ItemInfo(cfg.BaseURL(), id)
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
				err := ItemRefreshThumbnail(cfg.BaseURL(), id)
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
				thumbnail, err := ItemThumbnail(cfg.BaseURL(), itemId)
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
		opts, defaultErr := ItemUpdateOptions{}.WithDefaults()
		cmd := &cobra.Command{
			Use:   "update",
			Short: "Update item",
			RunE: func(cmd *cobra.Command, args []string) error {
				if defaultErr != nil {
					return defaultErr // should be good enough
				}
				resp, err := ItemUpdate(cfg.BaseURL(), opts)
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
