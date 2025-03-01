package api

// wrapper for every path in the eagle library.

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"

	"github.com/eissar/nest/api/endpoints"
	"github.com/labstack/echo/v4"
)

type EagleApiResponse struct {
	Status string
	Data   []interface{} // optional
}

//	type EagleResponse struct {
//		Status string
//		Data   []Item // optional
//	}

// for endpoints that return an array of data.
type EagleData struct {
	Status string `json:"status"`
}

// maybe
func (data EagleData) GetData() {}

type EagleMessage struct {
	Status string
	Data   any
}

type EagleApiErr struct {
	Message  string
	Endpoint endpoints.Endpoint
	Err      error
}

func (e *EagleApiErr) Error() string {
	return fmt.Sprintf("eagle api error calling path=%s docurl=%s err=%v ", e.Endpoint.Path, e.Endpoint.HelpUri(), e.Err)
}

type ApiKeyErr struct {
	Message string
}

func (e *ApiKeyErr) Error() string {
	return fmt.Sprintf("eagleapi: api key invalid; err=%s", e.Message)
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
func InvokeEagleAPI(req *http.Request, body interface{}) (*EagleApiResponse, error) {
	var result EagleApiResponse
	key, err := getApiKey()
	if err != nil {
		return nil, err
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

// mutates r
func addTokenAndEncodeQueryParams(r *http.Request) error {
	key, err := getApiKey()
	if err != nil {
		return err
	}

	query := r.URL.Query()
	query.Add("token", key)
	r.URL.RawQuery = query.Encode()
	return nil
}

// all responses have a status
// (excl. /api/library/icon)
func InvokeEagleAPIV1(req *http.Request) (result *EagleData, e error) {
	err := addTokenAndEncodeQueryParams(req)
	if err != nil {
		return result, err
	}

	// make the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return result, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	// parse the response
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return result, fmt.Errorf("error decoding response: %v", err)
	}

	if result.Status != "success" {
		return result, fmt.Errorf("error decoding response: result object's response was not `success`, but instead, %s ", result.Status)
	}

	return result, nil
}

type Ptr[T any] interface{ *T }

// all responses have a status
// (excl. /api/library/icon)
// populates pointer with response from req
func InvokeEagleAPIV2[T any](req *http.Request, v *T) error {
	if v == nil {
		return fmt.Errorf("v cannot be nil.")
	}
	err := addTokenAndEncodeQueryParams(req)
	if err != nil {
		return err
	}

	// make the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode > 299 || resp.StatusCode < 200 {
		var error_message any
		err = json.NewDecoder(resp.Body).Decode(&error_message)
		if err != nil {
			return fmt.Errorf("error decoding response: %w", err)
		}
		return fmt.Errorf("response code from eagle was not 2XX: response: %v", error_message)
	}

	// parse the response
	err = json.NewDecoder(resp.Body).Decode(&v)
	if err != nil {
		return fmt.Errorf("error decoding response: %v", err)
	}
	return nil
}

func wrapperHandler(c echo.Context) error {
	if c.Request().Method == "GET" {

	}
	return c.String(200, c.Request().URL.Path)
}

func RegisterEagleWrapper(g *echo.Group) {
	//for _, ep := range endpoints.Application {
	//for _, ep := range endpoints.Folder {
	for _, ep := range endpoints.Item {
		if ep.Method == "GET" {
			g.GET(ep.Path, wrapperHandler)
		} else if ep.Method == "POST" {
			g.POST(ep.Path, wrapperHandler)
		}
	}
	//for _, ep := range endpoints.Library {
}

const (
	MaxEagleItemIDLength = 15
	eagleItemIDPattern   = `^[a-zA-Z0-9]+$` // Pre-compiled regular expression
)

var eagleItemIDRegex = regexp.MustCompile(eagleItemIDPattern)

// todo remove regex
func IsValidItemID(id string) bool {
	if len(id) >= MaxEagleItemIDLength {
		return false
	}
	return eagleItemIDRegex.MatchString(string(id))
}

// docs: https://api.eagle.cool/item/add-from-url

func RegisterGroupRoutes(g *echo.Group) {
	g.GET("*", wrapperHandler)
	//g.GET("/item/addFromURL", handleAddItemFromUrl)
}

func RegisterRootRoutes(server *echo.Echo) {
	server.GET("/http\\:*", func(c echo.Context) error {
		return c.String(200, c.Request().URL.Path)
	})
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
