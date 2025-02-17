package api

// wrapper for every path in the eagle library.

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
)

type Item struct {
	URL              string            `json:"url"`
	Name             string            `json:"name"`
	Website          string            `json:"website,omitempty"`
	Tags             []string          `json:"tags,omitempty"`
	Star             int               `json:"star,omitempty"`
	Annotation       string            `json:"annotation,omitempty"`
	ModificationTime int64             `json:"modificationTime,omitempty"`
	FolderID         string            `json:"folderId,omitempty"`
	Headers          map[string]string `json:"headers,omitempty"`
}

type Endpoint struct {
	Path    string
	HelpUri string
	Method  string
}

type EagleApiErr struct {
	Message  string
	Endpoint Endpoint
	Err      error
}

type ApiKeyErr struct {
	Message string
}

func (e *ApiKeyErr) Error() string {
	return fmt.Sprintf("eagleapi: api key invalid; err=%s", e.Message)
}

func (e *EagleApiErr) Error() string {
	return fmt.Sprintf("eagle api error calling path=%s docurl=%s err=%v ", e.Endpoint.Path, e.Endpoint.HelpUri, e.Err)
}

func getApiKey() (string, error) {
	accessToken := os.Getenv("EAGLE_API_KEY") // Get token from environment variable
	if accessToken == "" {
		return "", &ApiKeyErr{
			Message: "environment variable `EAGLE_API_KEY` is not set or is empty.",
		}
	}
	return accessToken, nil
}

func InvokeRaindropAPI(endpoint Endpoint, body interface{}) (string, error) {
	key, err := getApiKey()
	if err != nil {
		return "", err
	}
	fmt.Println(key)
	return "a", nil

	/*
		headers := map[string]string{
			"Authorization": "Bearer " + key,
			"Content-Type":  "application/json",
		}

		var requestBody []byte
		if body != nil {
			if bodyStr, ok := body.(string); ok {
				requestBody = []byte(bodyStr)
			} else {
				jsonBody, err := json.Marshal(body)
				if err != nil {
					return "", fmt.Errorf("error marshaling body to JSON: %v", err)
				}
				requestBody = jsonBody
			}
		}

		req, err := http.NewRequest(method, endpoint, bytes.NewBuffer(requestBody))
		if err != nil {
			return "", fmt.Errorf("error creating request: %v", err)
		}

		for k, v := range headers {
			req.Header.Set(k, v)
		}

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return "", fmt.Errorf("error making request: %v", err)
		}
		defer resp.Body.Close()

		responseBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", fmt.Errorf("error reading response body: %v", err)
		}

		if resp.StatusCode != http.StatusOK {
			return string(responseBody), fmt.Errorf("API Error: Status Code %d, Response: %s", resp.StatusCode, string(responseBody))
		}

		return string(responseBody), nil
	*/
}

// docs: https://api.eagle.cool/item/add-from-url
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

	if item.Headers != nil {
		for k, v := range item.Headers {
			req.Header.Set(k, v)
		}
	}

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

func RegisterGroupRoutes(g *echo.Group) {
	g.GET("/item/addFromURL", handleAddItemFromUrl)

}
