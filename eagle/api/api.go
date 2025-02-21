package api

// wrapper for every path in the eagle library.

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/eissar/nest/eagle/api/endpoints"
	"github.com/labstack/echo/v4"
)

type EagleApiErr struct {
	Message  string
	Endpoint endpoints.Endpoint
	Err      error
}

type EagleResponse struct {
	Status string
	Data   []Item
}

type ApiKeyErr struct {
	Message string
}

func (e *ApiKeyErr) Error() string {
	return fmt.Sprintf("eagleapi: api key invalid; err=%s", e.Message)
}

func (e *EagleApiErr) Error() string {
	return fmt.Sprintf("eagle api error calling path=%s docurl=%s err=%v ", e.Endpoint.Path, e.Endpoint.HelpUri(), e.Err)
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

// adds api key to request and
// returns *EagleResponse.
//
// TODO: we check if status == success anyways, so
// should we just return EagleResponse.Data?
func InvokeRaindropAPI(req *http.Request, body interface{}) (*EagleResponse, error) {
	var result EagleResponse
	key, err := getApiKey()
	if err != nil {
		return &result, err
	}

	query := req.URL.Query()
	query.Add("key", key)

	req.URL.RawQuery = query.Encode()

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return &result, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return &result, fmt.Errorf("error decoding response: %v", err)
	}

	// TODO: do all responses have a status key?
	if result.Status != "success" {
		return &result, fmt.Errorf("error decoding response: result object's response was not `success`, but instead, %s ", result.Status)
	}

	return &result, nil
}

// docs: https://api.eagle.cool/item/add-from-url

func RegisterGroupRoutes(g *echo.Group) {
	g.GET("/item/addFromURL", handleAddItemFromUrl)

}

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
