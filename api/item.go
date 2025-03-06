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

// site, annotation, tags, folderid
type ItemProto struct {
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
type ListItem struct {
	Id               string   `json:"id"`
	URL              string   `json:"url"`
	Name             string   `json:"name"`
	Website          string   `json:"website"`
	Tags             []string `json:"tags"`
	Annotation       string   `json:"annotation"`
	ModificationTime int64    `json:"modificationTime"`
	//FolderID         string   `json:"folderId"`
}

type ItemUrl struct {
}

// no folder Id
type BulkItem struct {
	Item
	//FolderId string `json:"omitempty`
}

//FolderID         string   `json:"folderId,omitempty"`

func AddItemFromURL(baseURL string, item Item) (map[string]interface{}, error) {
	itemJSON, err := json.Marshal(item)
	if err != nil {
		return nil, fmt.Errorf("error marshaling item to JSON: %v", err)
	}

	req, err := http.NewRequest("POST", baseURL+"/api/item/addFromURL", bytes.NewBuffer(itemJSON))
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

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	if resp.StatusCode != http.StatusOK { // Check for non-200 status codes
		return nil, fmt.Errorf("server returned non-200 status: %d, Response: %v", resp.StatusCode, result)
	}

	return result, nil
}

type ListData struct {
	EagleData
	Data []any `json:"data"`
}

type Options struct {
	Limit int `json:"limit,omitempty"`
}

// creates an *http.Request and sends to InvokeEagleAPIV1
func ListV2(baseUrl string, limit int) (*ListData, error) {
	/*
		PARAMS
			limit
			The number of items to be displayed. the default number is 200
			offset
			Offset a collection of results from the api. Start with 0.
			orderBy
			The sorting order.CREATEDATE , FILESIZE , NAME , RESOLUTION , add a minus sign for descending order: -FILESIZE
			keyword
			Filter by the keyword
			ext
			Filter by the extension type, e.g.: jpg ,  png
			tags
			Filter by tags. Use , to divide different tags. E.g.: Design, Poster
			folders
			Filter by Folders.  Use , to divide folder IDs. E.g.: KAY6NTU6UYI5Q,KBJ8Z60O88VMG
	*/
	ep, ok := endpoints.Item["list"]
	if !ok {
		return nil, fmt.Errorf("could not find endpoint `list` in endpoints.")
	}
	// TODO: validate parameters?

	uri := baseUrl + ep.Path

	req, err := http.NewRequest(ep.Method, uri, nil) // method, url, body
	if err != nil {
		return nil, fmt.Errorf("list: error creating request err=%w", err)
	}

	// add query params
	query := req.URL.Query()
	if limit > 0 {
		query.Add("limit", strconv.Itoa(limit))
	}
	req.URL.RawQuery = query.Encode()
	// fmt.Println("query here:", req.URL.Query().Encode())

	var a *ListData
	err = InvokeEagleAPIV2(req, &a)
	if err != nil {
		return nil, err
	}
	return a, nil
}

type ListDataV3 struct {
	EagleData
	Data []ListItem
}

// creates an *http.Request and sends to InvokeEagleAPIV1
func ListV3(baseUrl string, limit int) (*ListDataV3, error) {
	/*
		PARAMS
			limit
			The number of items to be displayed. the default number is 200
			offset
			Offset a collection of results from the api. Start with 0.
			orderBy
			The sorting order.CREATEDATE , FILESIZE , NAME , RESOLUTION , add a minus sign for descending order: -FILESIZE
			keyword
			Filter by the keyword
			ext
			Filter by the extension type, e.g.: jpg ,  png
			tags
			Filter by tags. Use , to divide different tags. E.g.: Design, Poster
			folders
			Filter by Folders.  Use , to divide folder IDs. E.g.: KAY6NTU6UYI5Q,KBJ8Z60O88VMG
	*/
	ep, ok := endpoints.Item["list"]
	if !ok {
		return nil, fmt.Errorf("could not find endpoint `list` in endpoints.")
	}
	// TODO: validate parameters?

	uri := baseUrl + ep.Path

	req, err := http.NewRequest(ep.Method, uri, nil) // method, url, body
	if err != nil {
		return nil, fmt.Errorf("list: error creating request err=%w", err)
	}

	// add query params
	query := req.URL.Query()
	if limit > 0 {
		query.Add("limit", strconv.Itoa(limit))
	}
	req.URL.RawQuery = query.Encode()
	// fmt.Println("query here:", req.URL.Query().Encode())

	var a *ListDataV3
	err = InvokeEagleAPIV2(req, &a)
	if err != nil {
		return nil, err
	}
	return a, nil
}

type itemFromPath struct {
	Path       string   `json:"path"`                 // Required, the path of the local file.
	Name       string   `json:"name,omitempty"`       // Required, the name of the image to be added. (not really req)
	Website    string   `json:"website,omitempty"`    // The Address of the source of the image.
	Annotation string   `json:"annotation,omitempty"` // The annotation for the image.
	Tags       []string `json:"tags,omitempty"`       // Tags for the image.
	FolderId   string   `json:"folderId,omitempty"`   // If this parameter is defined, the image will be added to the corresponding folder.
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

func NewItemFromPath(filePath string) (obj itemFromPath, err error) {
	err = obj.setPath(filePath)
	return obj, err
}

// returns status only.
// use ConstructItemFromPath to build args
func AddItemFromPath(baseURL string, item itemFromPath) error {
	ep, ok := endpoints.Item["addFromPath"]
	if !ok {
		return fmt.Errorf("could not find endpoint `addFromPath` in endpoints.")
	}
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

type ThumbnailData struct {
	Status        string `json:"status"`
	ThumbnailPath string `json:"data"`
}

// creates an *http.Request and sends to InvokeEagleAPIV1
func Thumbnail(baseUrl string, itemId string) (s string, err error) {
	ep, ok := endpoints.Item["thumbnail"]
	if !ok {
		return s, fmt.Errorf("could not find endpoint `list` in endpoints.")
	}

	uri := baseUrl + ep.Path

	req, err := http.NewRequest(ep.Method, uri, nil) // method, url, body
	if err != nil {
		return s, fmt.Errorf("list: error creating request err=%w", err)
	}

	// add query params
	query := req.URL.Query()

	if !IsValidItemID(itemId) {
		return s, fmt.Errorf("list: error creating request err= itemId parameter malformed or empty.")
	}
	query.Add("id", itemId)
	req.URL.RawQuery = query.Encode()

	var a ThumbnailData
	err = InvokeEagleAPIV2(req, &a)
	if err != nil {
		return s, err
	}

	if escapedPath, err := url.PathUnescape(a.ThumbnailPath); err != nil {
		return a.ThumbnailPath, fmt.Errorf("could not url decode path response from eagle server err=%w", err)
	} else {
		return escapedPath, nil
	}
}
