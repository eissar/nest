package nest

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
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

type ThumbnailData struct {
	Status        string `json:"status"`
	ThumbnailPath string `json:"data"`
}

func GetEagleThumbnail(cfg *config.NestConfig, itemId string) (string, error) {
	thumbnail, err := api.Thumbnail(cfg.BaseURL(), itemId)
	if err != nil {
		return "", fmt.Errorf("getEagleThumbnail: err=%w", err)
	}
	return thumbnail, nil
}

var allowed_filetypes = []string{".jpeg", ".jpg", ".png", ".gif", ".svg", ".webp", ".avif"}

// tries to find the actual from the response
// of request api/item/thumbnail. checks if there are
// any files matching `allowed_filetypes`.
// also calls `url.PathUnescape` on the url.
func resolveThumbnailPath(t string) (string, error) {
	var thumbnail string // output

	if !strings.HasSuffix(t, "_thumbnail.png") {
		// should already the full-resolution file.
		thumbnail = t
	} else {
		thumbnailRoot := strings.TrimSuffix(thumbnail, "_thumbnail.png")

		for _, typ := range allowed_filetypes {
			joinedPath := thumbnailRoot + typ
			if _, err := os.Stat(joinedPath); err == nil {
				thumbnail = joinedPath // if any exists
				break
			}
		}
	}

	thumbnail, err := url.PathUnescape(t)
	if err != nil {
		return thumbnail, fmt.Errorf("resolvethumb: error cleaning thumbnail path: %s:", err.Error())
	}
	return thumbnail, nil
}

// on my device thumbnail ONLY end with _thumbnail.png or they do not exist.
// this returns the full file path
func GetEagleThumbnailFullRes(cfg *config.NestConfig, itemId string) (string, error) {
	thumbnail, err := GetEagleThumbnail(cfg, itemId)
	if err != nil {
		return thumbnail, fmt.Errorf("getEagleThumbnail: err=%w", err)
	}

	thumbnail, err = resolveThumbnailPath(thumbnail)
	if err != nil {
		return thumbnail, fmt.Errorf("getEagleThumbnail: err=%w", err)
	}

	thumbnail, err = url.PathUnescape(thumbnail)
	if err != nil {
		return thumbnail, fmt.Errorf("error cleaning thumbnail path: %s", err.Error())
	}

	//  TODO: fallback list all files other than metadata.json & _thumbnail.png?
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
