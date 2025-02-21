package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
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
