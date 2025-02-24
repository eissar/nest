package nest

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/eissar/nest/config"
	"github.com/eissar/nest/eagle/api"
	"github.com/eissar/nest/eagle/api/endpoints"
)

type nestCfg config.NestConfig

type EagleServerInfo struct {
	Status string `json:"status"`
	Data   struct {
		Version string `json:"version"`
	} `json:"data"`
}
type UrlBuilder struct {
	Url   *url.URL
	Query *url.Values
}
type EagleItemId string

const (
	MaxEagleItemIDLength = 15
	eagleItemIDPattern   = `^[a-zA-Z0-9]+$` // Pre-compiled regular expression
)

var eagleItemIDRegex = regexp.MustCompile(eagleItemIDPattern)

func (id EagleItemId) IsValid() bool {
	if len(id) >= MaxEagleItemIDLength {
		return false
	}
	return eagleItemIDRegex.MatchString(string(id))
}
func IsValidEagleItemId(itemId EagleItemId) bool {
	if itemId.IsValid() {
		return true
	} else {
		return false
	}
}

// validateIsEagleServerRunning checks if the Eagle server is running at the specified URL.
func validateIsEagleServerRunning(url string) (bool, error) {
	client := &http.Client{
		Timeout: 5 * time.Second, // Add a timeout to prevent indefinite hangs
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return false, fmt.Errorf("error creating request: %w", err)
	}
	//It is a good practice to set the header
	//req.Header.Set("User-Agent", "EagleServerValidator")

	resp, err := client.Do(req)
	if err != nil {
		return false, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close() // Ensure the response body is closed

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("eagle API request failed with status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, fmt.Errorf("error reading response body: %w", err)
	}

	var data EagleServerInfo
	if err := json.Unmarshal(body, &data); err != nil {
		return false, fmt.Errorf("error unmarshalling JSON response: %w", err)
	}

	// Check for specific response data (optional, but good for robustness)
	if data.Status == "success" && data.Data.Version != "" {
		return true, nil // Eagle server is running and data is valid
	}

	return false, nil
}

func (u UrlBuilder) String() string {
	u.Url.RawQuery = u.Query.Encode()
	return u.Url.String()
}
func Uri(cfg *config.NestConfig, path string) *UrlBuilder {
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	h := fmt.Sprintf("%s:%s", cfg.Host, strconv.Itoa(cfg.Port))
	uri := &url.URL{
		Scheme: "http",
		Host:   h,
		Path:   path,
	}

	queryValues := url.Values{}
	queryValues.Add("key", cfg.ApiKey)
	return &UrlBuilder{
		Url:   uri,
		Query: &queryValues,
	}
}

type ItemIds []string

func GetEagleThumbnailFromId(cfg *config.NestConfig, id string) (map[string]interface{}, error) {
	client := &http.Client{
		Timeout: 5 * time.Second, // Add a timeout to prevent indefinite hangs
	}
	ep := endpoints.Item["thumbnail"]

	builder := Uri(cfg, ep.Path)
	builder.Query.Add("id", id)

	req, err := http.NewRequest("GET", builder.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("geteagleThumbnailFromId: error while creating new request err=%s", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("geteagleThumbnailFromId: error making request: %w", err)
	}
	defer resp.Body.Close() // Ensure the response body is closed

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("eagle API request failed with status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}
	var a map[string]interface{}
	err = json.Unmarshal(body, &a)
	//fmt.Printf("body: %v\n", a["body"])
	if err != nil {
		log.Fatalf("getEagleThumbnailFromId: %s", err.Error())
	}
	return a, nil
}

var allowed_filetypes = []string{"JPEG", "PNG", "GIF", "SVG", "WebP", "AVIF"}

func GetEaglePathFromId(cfg *config.NestConfig, id string) (string, error) {
	client := &http.Client{
		Timeout: 5 * time.Second, // Add a timeout to prevent indefinite hangs
	}
	ep := endpoints.Item["thumbnail"]

	builder := Uri(cfg, ep.Path)
	builder.Query.Add("id", id)

	req, err := http.NewRequest(ep.Method, builder.String(), nil)
	if err != nil {
		return "", fmt.Errorf("getEagleThumbnailFromId: error while creating new request err=%s", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("getEagleThumbnailFromId: error making request: %w", err)
	}
	defer resp.Body.Close() // Ensure the response body is closed

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("eagle API request failed with status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %w", err)
	}
	var responseData EagleDataMessage
	err = json.Unmarshal(body, &responseData)
	//fmt.Printf("body: %v\n", a["body"])
	if err != nil {
		log.Fatalf("getEagleThumbnailFromId: %s", err.Error())
	}
	return responseData.Data, nil
}

type ThumbnailData struct {
	Status        string `json:"status"`
	ThumbnailPath string `json:"data"`
}

func GetEagleThumbnail(cfg *config.NestConfig, itemId EagleItemId) (ThumbnailData, error) {
	client := &http.Client{
		Timeout: 5 * time.Second, // Add a timeout to prevent indefinite hangs
	}
	var out ThumbnailData
	ep := endpoints.Item["thumbnail"]

	builder := Uri(cfg, ep.Path)
	builder.Query.Add("id", string(itemId))

	req, err := http.NewRequest(ep.Method, builder.String(), nil)
	if err != nil {
		return out, fmt.Errorf("geteagleThumbnailFromId: error while creating new request err=%w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return out, fmt.Errorf("geteagleThumbnailFromId: error making request: %w", err)
	}
	defer resp.Body.Close() // Ensure the response body is closed

	if resp.StatusCode != http.StatusOK {
		return out, fmt.Errorf("eagle API request failed with status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return out, fmt.Errorf("error reading response body: err=%w", err)
	}
	err = json.Unmarshal(body, &out)
	if err != nil {
		return out, fmt.Errorf("getEagleThumbnailFromId: error while unmarshalling data."+
			"has schema changed? check docs for endpoint=`%s` "+
			"<https://api.eagle.cool/item/thumbnail> err=%w ", ep.Path, err)
	}

	out.ThumbnailPath, err = url.PathUnescape(out.ThumbnailPath)
	if err != nil {
		return out, fmt.Errorf("getEagleThumbnailFromId: error while parsing data."+
			"err=%w", err)
	}

	if out.Status != "success" {
		return out, fmt.Errorf("getEagleThumbnailFromId: error while recieving data: "+
			"docs for endpoint=`%s` <https://api.eagle.cool/item/thumbnail> err= "+
			"resp.Status is `%s` and not `success`.", ep.Path, out.Status)
	}

	return out, nil
}

// on my device thumbnail ONLY end with _thumbnail.png or they do not exist.
// this returns the full file path if there is no thumbnail.
func GetEagleThumbnailV2(cfg *config.NestConfig, itemId string) (string, error) {
	baseUrl := "http://" + cfg.Host + ":" + strconv.Itoa(cfg.Port)
	thumbnail, err := api.Thumbnail(baseUrl, itemId)
	if err != nil {
		return "", fmt.Errorf("getEagleThumbnail: err=%w", err)
	}
	return thumbnail, nil
}

func GetListV0(cfg *config.NestConfig) (any, error) {
	baseUrl := "http://" + cfg.Host + ":" + strconv.Itoa(cfg.Port)
	_, err := api.ListV2(baseUrl, 5)
	if err != nil {
		log.Fatal(err.Error())
	}

	return nil, nil
}

/*
func getEagleThumbnailsFromIds() {
	client := &http.Client{
		Timeout: 5 * time.Second, // Add a timeout to prevent indefinite hangs
	}
	cfg := get
	endpoint := buildUri(,)
	req, err := http.NewRequest()
}
*/
