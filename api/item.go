package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	_ "net/url"
	"os"
	"path/filepath"
	"strconv"

	"github.com/eissar/nest/api/endpoints"
)

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
type itemFromPath struct {
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

// resolves p string and sets path.
func (i *itemFromPath) setPath(p string) error {
	file_path, err := filepath.Abs(p)
	if err != nil {
		return fmt.Errorf("[ERROR] could not resolve %s err=%w\n", p, err)
	}
	file_path = filepath.ToSlash(file_path)

	_, err = os.Stat(file_path)
	if err != nil {
		return fmt.Errorf("[ERROR] could not resolve %s err=%w\n", p, err)
	}

	i.Path = file_path
	return nil
}

// construct private api.Item type from just filepath
func NewItemFromPath(filePath string) (obj itemFromPath, err error) {
	err = obj.setPath(filePath)
	return obj, err
}

// start api endpoints

// use cleaner abstraction `api.Request` like in ItemInfo

//- [X] /api/item/addFromURL
//- [ ] /api/item/addFromURLs
//- [ ] /api/item/addFromPath
//- [ ] /api/item/addFromPaths
//- [ ] /api/item/addBookmark
//- [X] /api/item/info
//- [X] /api/item/thumbnail
//- [X] /api/item/list
//- [ ] /api/item/moveToTrash
//- [ ] /api/item/refreshPalette
//- [ ] /api/item/refreshThumbnail
//- [X] /api/item/update

// old; deprecated
func AddItemFromURL(baseUrl string, item Item) (map[string]any, error) {
	// convert to json
	itemJSON, err := json.Marshal(item)
	if err != nil {
		return nil, fmt.Errorf("error marshaling item to JSON: %v", err)
	}

	req, err := http.NewRequest("POST", baseUrl+"/api/item/addFromURL", bytes.NewBuffer(itemJSON))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	//if item.Headers != nil {
	//	for k, v := range item.Headers {
	//		req.Header.Set(k, v)
	//	}
	//}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	var result map[string]any
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	if resp.StatusCode != http.StatusOK { // Check for non-200 status codes
		return nil, fmt.Errorf("server returned non-200 status: %d, Response: %v", resp.StatusCode, result)
	}

	return result, nil
}

// only returns status
func ItemAddFromUrl(baseUrl string) error {
	ep := endpoints.ItemAddFromURL
	uri := baseUrl + ep.Path

	req, err := http.NewRequest(ep.Method, uri, nil) // method, url, body
	if err != nil {
		return fmt.Errorf("list: error creating request err=%w", err)
	}

	var resp EagleResponse
	err = InvokeEagleAPIV2(req, &resp)
	if err != nil {
		return fmt.Errorf("error invoking eagle api err=%v", err)
	}

	if resp.Status != "success" {
		return fmt.Errorf("response status recieved from eagle was not `success` message=%v", resp.Status)
	}

	return nil
}

func ItemAddBookmark(baseUrl string) error {

	return nil
}

// creates an *http.Request and sends to InvokeEagleAPIV1
func ItemList(baseUrl string, opts ItemListOptions) ([]ListItem, error) {
	ep := endpoints.ItemList

	uri := baseUrl + ep.Path

	req, err := http.NewRequest(ep.Method, uri, nil) // method, url, body
	if err != nil {
		return nil, fmt.Errorf("list: error creating request err=%w", err)
	}

	// add query params
	query := req.URL.Query()
	if opts.Limit > 0 {
		query.Add("limit", strconv.Itoa(opts.Limit))
	}
	req.URL.RawQuery = query.Encode()
	// fmt.Println("query here:", req.URL.Query().Encode())

	// TODO: validate parameters

	var resp struct {
		EagleData
		Data []ListItem `json:"data"`
	}
	err = InvokeEagleAPIV2(req, &resp)
	if err != nil {
		return nil, err
	}

	// if a.Status != "success" ...

	return resp.Data, nil
}

// returns status only.
// use NewItemFromPath to build args
func AddItemFromPath(baseURL string, item itemFromPath) error {
	ep := endpoints.ItemAddFromPath

	uri := baseURL + ep.Path
	body := fmt.Appendf(nil, `{"path": "%s"}`, item.Path)

	req, err := http.NewRequest(ep.Method, uri, bytes.NewReader(body)) // method, url, body
	if err != nil {
		return fmt.Errorf("list: error creating request err=%w", err)
	}

	var a *EagleMessage
	err = InvokeEagleAPIV2(req, &a)
	if err != nil {
		return err
	}

	if a.Status != "success" {
		return fmt.Errorf("response status recieved from eagle was not `success` message=%v", a)
	}

	fmt.Println(a)

	return nil
}

// deprecated
func ItemInfoV0(baseUrl string, id string) (respItem ApiItem, err error) {
	ep := endpoints.ItemList

	uri := baseUrl + ep.Path

	req, err := http.NewRequest(ep.Method, uri, nil) // method, url, body
	if err != nil {
		return respItem, fmt.Errorf("list: error creating request err=%w", err)
	}

	// validate id
	if !IsValidItemID(id) {
		return respItem, fmt.Errorf("list: error creating request err= itemId parameter malformed or empty.")
	}

	// set query params
	query := req.URL.Query()
	query.Add("id", id)
	req.URL.RawQuery = query.Encode()

	var resp struct {
		EagleResponse
		Data ApiItem `json:"data"`
	}
	err = InvokeEagleAPIV2(req, &resp)
	if err != nil {
		return respItem, err
	}

	// if a.Status != "success" ...

	return respItem, nil
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

// * id Required, the ID of the item to be modified
// * tags Optional, tags
// * annotation Optional, annotations
// * url Optional, the source url
// * star Optional, ratings
// pointers represent optional keys and null represents unset
type ItemUpdateOptions struct {
	ID         string    `json:"id"`
	Tags       *[]string `json:"tags"`
	Annotation *string   `json:"annotation"`
	URL        *string   `json:"url"`
	Star       *int      `json:"star"`
}

func ItemUpdate(baseUrl string, opts ItemUpdateOptions) (respItem ApiItem, err error) {
	ep := endpoints.ItemUpdate

	uri := baseUrl + ep.Path

	req, err := http.NewRequest(ep.Method, uri, nil) // method, url, body
	if err != nil {
		return respItem, fmt.Errorf("list: error creating request err=%w", err)
	}

	// validate id
	if !IsValidItemID(opts.ID) {
		return respItem, fmt.Errorf("list: error creating request err= itemId parameter malformed or empty.")
	}

	// set query params
	query := req.URL.Query()
	query.Add("id", opts.ID)
	req.URL.RawQuery = query.Encode()

	return respItem, nil
}

// returns thumbnail path and error
func ItemThumbnail(baseUrl string, itemId string) (string, error) {
	ep := endpoints.ItemThumbnail

	uri := baseUrl + ep.Path

	req, err := http.NewRequest(ep.Method, uri, nil) // method, url, body
	if err != nil {
		return "", fmt.Errorf("list: error creating request err=%w", err)
	}

	// validate query params
	if !IsValidItemID(itemId) {
		return "", fmt.Errorf("list: error creating request err= itemId parameter malformed or empty.")
	}

	// add query params
	query := req.URL.Query()

	query.Add("id", itemId)
	req.URL.RawQuery = query.Encode()

	var a ThumbnailData
	err = InvokeEagleAPIV2(req, &a)
	if err != nil {
		return "", err
	}

	if escapedPath, err := url.PathUnescape(a.ThumbnailPath); err != nil {
		return a.ThumbnailPath, fmt.Errorf("could not url decode path response from eagle server err=%w", err)
	} else {
		return escapedPath, nil
	}
}
