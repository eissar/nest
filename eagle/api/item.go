package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/eissar/nest/eagle/api/endpoints"
)

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

func List(baseURL string, limit int) (EagleResponse, error) {
	var result EagleResponse

	req, err := http.NewRequest(http.MethodGet, baseURL+"/api/item/list", http.NoBody)
	if err != nil {
		return result, fmt.Errorf("error initializing request: %v", err)
	}
	query := req.URL.Query()

	query.Add("limit", strconv.Itoa(limit))
	req.URL.RawQuery = query.Encode()
	//fmt.Printf("url: %v\n", req.URL.RequestURI())

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return result, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return result, fmt.Errorf("error decoding response: %v", err)
	}

	if result.Status != "success" {
		fmt.Errorf("error decoding response: result object's response was not `success`, but instead, %s ", result.Status)
	}

	resp.Body.Close()
	return result, nil
}

func ListV1(baseURL string, limit int) (EagleResponse, error) {
	var result EagleResponse

	req, err := http.NewRequest(http.MethodGet, baseURL+"/api/item/list", http.NoBody)
	if err != nil {
		return result, fmt.Errorf("error initializing request: %v", err)
	}
	query := req.URL.Query()
	query.Add("limit", strconv.Itoa(limit))
	//fmt.Printf("url: %v\n", req.URL.RequestURI())
	req.Header.Set("Content-Type", "application/json")

	InvokeEagleAPIV1(req)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return result, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return result, fmt.Errorf("error decoding response: %v", err)
	}

	if result.Status != "success" {
		return result, fmt.Errorf("error decoding response: result object's response was not `success`, but instead, %s ", result.Status)
	}

	resp.Body.Close()
	return result, nil
}

// creates an *http.Request and sends to InvokeEagleAPIV1
func ListV2(baseUrl string) (*EagleData, error) {
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

	return InvokeEagleAPIV1(req)
}
