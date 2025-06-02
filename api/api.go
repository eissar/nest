package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"syscall"

	"github.com/eissar/nest/api/endpoints"
	"github.com/labstack/echo/v4"
	"golang.org/x/sys/windows"
)

var (
	ErrStatusErr = fmt.Errorf("response key 'status' was not 'success'")
)

// #region errors

type EagleApiErr struct {
	Message  string
	Endpoint endpoints.Endpoint
	Err      error
}

func (e *EagleApiErr) Error() string {
	return fmt.Sprintf("eagle api error calling path=%s docurl=%s err=%v ", e.Endpoint.Path, e.Endpoint.HelpUri(), e.Err)
}

var EagleNotOpenErr = fmt.Errorf("Eagle is not open.")

// #endregion errors

// #region types

type ApiKeyErr struct {
	Message string
}

func (e *ApiKeyErr) Error() string {
	return fmt.Sprintf("eagleapi: api key invalid; err=%s", e.Message)
}

type EagleApiResponse struct {
	Status string
	Data   []interface{} // optional
}

// maybe? func (data EagleResponse) GetData() {}
type EagleResponse struct {
	Status string `json:"status"`
}

type EagleMessage struct {
	EagleResponse
	Data any
}

// for endpoints that return an array of strings.
type EagleArray struct {
	EagleResponse
	Data []string
}

// #endregion types

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

// populates v with response from req
func Request[T any](method string, url string, body io.Reader, urlParam *url.Values, v *T) error {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return fmt.Errorf("Request: error creating request err=%w", err)
	}

	if urlParam != nil {
		req.URL.RawQuery = urlParam.Encode()
	}
	err = InvokeEagleAPIV2(req, &v)
	if err != nil {
		return err
	}

	return nil
}

func IsEagleNotOpenErr(err error) bool {
	// windows, linux
	if errors.Is(err, windows.WSAECONNREFUSED) || errors.Is(err, syscall.ECONNREFUSED) {
		return true
	}
	return false
	// var sysErr *os.SyscallError
	// if ok := errors.As(err, &sysErr); ok {
	// 	if errno, ok := sysErr.Err.(syscall.Errno); ok && int(errno) == 10061 {
	// 		// <https://github.com/ddev/ddev/blob/ec7870af3af6356cfe26282b6d64e551e1891544/pkg/netutil/netutil.go#L36>

	// 		return true
	// 	}
	// }
}

// all responses have a `status` field (excl. /api/library/icon)
// populates pointer v with response from req
func InvokeEagleAPIV2[T any](req *http.Request, v *T) error {
	if v == nil {
		return fmt.Errorf("v cannot be nil.")
	}

	var requestBodyBytes []byte // Store the body here
	var readErr error
	requestBodyBytes, readErr = io.ReadAll(req.Body)
	// It's crucial to close the original body after you've read it.
	req.Body.Close()
	if readErr != nil {
		// If reading fails, you won't have the body for the error message.
		// You could return an error here or log it and continue without the body string.
		// For this minimal example, we'll proceed, and requestBodyBytes might be empty/nil.
		// Consider returning: return fmt.Errorf("failed to read request body for error message: %w", readErr)
	}
	// IMPORTANT: Replace req.Body so client.Do can read it again.
	req.Body = io.NopCloser(bytes.NewBuffer(requestBodyBytes))

	err := addTokenAndEncodeQueryParams(req)
	if err != nil {
		return err
	}

	// make the request
	client := &http.Client{}
	resp, err := client.Do(req)

	if IsEagleNotOpenErr(err) {
		return EagleNotOpenErr
	}
	if err != nil {
		oErr, ok := err.(*net.OpError)
		if ok {
			//if IsEagleNotOpenErr(oErr) {
			//	fmt.Println("NOT OPEN")
			//	return EagleNotOpenErr
			//}

			return fmt.Errorf("InvokeEagleAPI: unknown err=%v", oErr)
		}
		return fmt.Errorf("InvokeEagleAPI: unknown err=%v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode > 299 || resp.StatusCode < 200 {
		var error_message any
		err = json.NewDecoder(resp.Body).Decode(&error_message)
		if err != nil {
			return fmt.Errorf("error decoding response: %w", err)
		}

		return fmt.Errorf("response code from eagle was not 2XX: response: %v; request body: %v", error_message, string(requestBodyBytes))
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

// #region routes

func RegisterGroupRoutes(g *echo.Group) {
	g.GET("*", wrapperHandler)
	//g.GET("/item/addFromURL", handleAddItemFromUrl)
}

func RegisterRootRoutes(server *echo.Echo) {
	server.GET("/http\\:*", func(c echo.Context) error {
		return c.String(200, c.Request().URL.Path)
	})
}

// #endregion routes
