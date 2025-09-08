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
	"reflect"
	"regexp"
	"strconv"
	"strings"
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

// sentinel errors
var LibraryIsAlreadyTargetErr = errors.New("Library is already active")
var EagleNotOpenOrUnavailableErr = fmt.Errorf("Eagle is not open or unavailable.")

// constructor
func GetCurrentLibraryIsAlreadyTargetError(currLib string) error {
	return fmt.Errorf("%w:  %s", LibraryIsAlreadyTargetErr, currLib)
}

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

// StructToURLValues converts a struct to a url.Values map based on its json tags.
// This allows you to easily serialize a struct into URL query parameters.
func StructToURLValues(data interface{}) (url.Values, error) {
	// The url.Values type is a map[string][]string, which is what http.Request.URL.Query() returns.
	// It's the standard way to represent query parameters in Go.
	values := url.Values{}

	// Use reflection to inspect the struct.
	// We expect data to be a struct, so we get its value.
	v := reflect.ValueOf(data)
	if v.Kind() == reflect.Ptr {
		// If it's a pointer, dereference it to get the struct.
		v = v.Elem()
	}

	// Ensure we are working with a struct.
	if v.Kind() != reflect.Struct {
		return nil, fmt.Errorf("StructToURLValues only accepts structs; got %T", data)
	}

	// Get the type of the struct to access its fields and tags.
	t := v.Type()

	// Iterate over all the fields of the struct.
	for i := 0; i < v.NumField(); i++ {
		fieldValue := v.Field(i)
		fieldType := t.Field(i)

		// Get the json tag for the current field.
		jsonTag := fieldType.Tag.Get("json")

		// Skip this field if the json tag is "-"
		if jsonTag == "-" {
			continue
		}

		// Parse the tag to get the parameter name and options like "omitempty".
		tagParts := strings.Split(jsonTag, ",")
		paramName := tagParts[0]

		// If the paramName is empty, it means the field is unexported or has no tag.
		// We use the field name as a fallback, but this is often not desired.
		// A better practice is to ensure all exported fields have tags.
		if paramName == "" {
			// Skip unexported fields.
			if !fieldType.IsExported() {
				continue
			}
			paramName = fieldType.Name
		}

		// Check for the "omitempty" option.
		hasOmitempty := false
		if len(tagParts) > 1 {
			for _, part := range tagParts[1:] {
				if part == "omitempty" {
					hasOmitempty = true
					break
				}
			}
		}

		// If "omitempty" is present and the field has its zero value, skip it.
		if hasOmitempty && fieldValue.IsZero() {
			continue
		}

		// Convert the field's value to a string.
		var paramValue string
		switch fieldValue.Kind() {
		case reflect.String:
			paramValue = fieldValue.String()
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			paramValue = strconv.FormatInt(fieldValue.Int(), 10)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			paramValue = strconv.FormatUint(fieldValue.Uint(), 10)
		case reflect.Float32, reflect.Float64:
			paramValue = strconv.FormatFloat(fieldValue.Float(), 'f', -1, 64)
		case reflect.Bool:
			paramValue = strconv.FormatBool(fieldValue.Bool())
		case reflect.Slice:
			// Handle slices by joining elements with a comma.
			// Another common approach is to add multiple parameters with the same name.
			// e.g., ?tags=go&tags=web
			// We'll demonstrate the comma-separated approach here.
			sliceVal := reflect.ValueOf(fieldValue.Interface())
			var elements []string
			for j := 0; j < sliceVal.Len(); j++ {
				elements = append(elements, fmt.Sprint(sliceVal.Index(j).Interface()))
			}
			paramValue = strings.Join(elements, ",")
		default:
			// For other types, you might need more complex logic.
			// For this example, we'll just skip them.
			continue
		}

		// Add the key-value pair to our url.Values map.
		values.Add(paramName, paramValue)
	}

	return values, nil
}

// TODO: remove regex?
func IsValidItemID(id string) bool {
	if len(id) >= MaxEagleItemIDLength {
		return false
	}
	// return eagleItemIDRegex.MatchString(string(id))
	return true
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

func IsEagleNotOpenOrUnavailable(err error) bool {
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
	if req.Body != nil {
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
	}

	err := addTokenAndEncodeQueryParams(req)
	if err != nil {
		return err
	}

	// make the request
	client := &http.Client{}
	resp, err := client.Do(req)

	// fmt.Printf("resp.StatusCode: %v\n", resp.StatusCode)
	// fmt.Printf("req.URL: %v\n", req.URL)

	if err != nil {
		if IsEagleNotOpenOrUnavailable(err) {
			return EagleNotOpenOrUnavailableErr
		}
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
			// Warning: failed to decode JSON response, attempting to read as string
			fmt.Println("[WARN] Could not decode error message.")
			bodyBytes, _ := io.ReadAll(resp.Body)
			error_message = string(bodyBytes)
		}

		return fmt.Errorf("response code from eagle was not 2XX: response: %v; request body: %v", error_message, string(requestBodyBytes))
	}

	contentLength, err := strconv.Atoi(resp.Header.Get("Content-Length"))
	if err != nil {
		contentLength = -1
	}
	if contentLength > 0 && contentLength < 1024 {
		// readall
		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		if err := json.Unmarshal(respBody, &v); err != nil {
			b := string(respBody)
			if strings.Contains(b, "Library does not exist.") {
				fmt.Println("[ERR] Eagle response: Library does not exist. \nCheck that the library path exists and is accessible (e.g., not an unavailable network drive or disconnected volume)")
				os.Exit(1)
			}

			fmt.Printf("error decoding response. resp: %s\n", string(respBody))
		}
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
