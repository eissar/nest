package api

// wrapper for every path in the eagle library.
// [X] - application
//		[X] - tests
// [ ] - folder
// [ ] - item
//		[ ] - tests
//		[ ] - parameters
// [X] - library
//		[X] - tests

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"regexp"

	"github.com/eissar/nest/api/endpoints"
	"github.com/labstack/echo/v4"
)

// errors
var (
	ErrStatusErr = fmt.Errorf("response key 'status' was not 'success'")
)

// types
type EagleApiResponse struct {
	Status string
	Data   []interface{} // optional
}
type EagleData struct {
	Status string `json:"status"`
}

// maybe
func (data EagleData) GetData() {}

// cleaned up way to do this

type EagleMessage struct {
	EagleData
	Data any
}

// for endpoints that return an array of strings.
type EagleArray struct {
	EagleData
	Data []string
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

// all responses have a `status` field (excl. /api/library/icon)
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

// populates v with response from req
func Request[T any](method string, url string, body io.Reader, urlParam *url.Values, v *T) error {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return fmt.Errorf("list: error creating request err=%w", err)
	}

	if urlParam != nil {
		req.URL.RawQuery = urlParam.Encode()
	}
	err = InvokeEagleAPIV2(req, &v)
	if err != nil {
		return fmt.Errorf("api.Request error making request err=%w", err)
	}

	return nil
}

// all responses have a `status` field (excl. /api/library/icon)
// populates pointer v with response from req
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

// eagle item id
const (
	MaxEagleItemIDLength = 15
	eagleItemIDPattern   = `^[a-zA-Z0-9]+$` // Pre-compiled regular expression
)

var eagleItemIDRegex = regexp.MustCompile(eagleItemIDPattern)

// TODO: remove regex?
func IsValidItemID(id string) bool {
	if len(id) >= MaxEagleItemIDLength {
		return false
	}
	return eagleItemIDRegex.MatchString(string(id))
}

func RegisterGroupRoutes(g *echo.Group) {
	g.GET("*", wrapperHandler)
	//g.GET("/item/addFromURL", handleAddItemFromUrl)
}

func RegisterRootRoutes(server *echo.Echo) {
	server.GET("/http\\:*", func(c echo.Context) error {
		return c.String(200, c.Request().URL.Path)
	})
}
