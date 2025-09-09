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

	"github.com/eissar/nest/api"
	"github.com/eissar/nest/config"
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

type ThumbnailData struct {
	Status        string `json:"status"`
	ThumbnailPath string `json:"data"`
}

func GetEagleThumbnail(cfg *config.NestConfig, itemId string) (string, error) {
	thumbnail, err := api.ItemThumbnail(cfg.BaseURL(), itemId)
	if err != nil {
		return "", fmt.Errorf("getEagleThumbnail: err=%w", err)
	}

	thumbnail, err = url.PathUnescape(thumbnail)
	if err != nil {
		return thumbnail, fmt.Errorf("getEagleThumbnail: error cleaning thumbnail path: %s:", err.Error())
	}

	if _, err = os.Stat(thumbnail); err != nil {
		return thumbnail, fmt.Errorf("getEagleThumbnail: err=%w", err)
	}

	return thumbnail, nil
}

var allowed_filetypes = []string{".jpeg", ".jpg", ".png", ".gif", ".svg", ".webp", ".avif"}

// tries to find the actual filepath from the response
// of request api/item/thumbnail.
// also calls `url.PathUnescape` on the url.
// then checks if there are
// any files matching `allowed_filetypes`.
func resolveThumbnailPath(thumbnail string) (string, error) {
	thumbnail, err := url.PathUnescape(thumbnail)
	if err != nil {
		return thumbnail, fmt.Errorf("resolvethumb: error cleaning thumbnail path: %s:", err.Error())
	}

	if !strings.HasSuffix(thumbnail, "_thumbnail.png") {
		// should already the full-resolution file.
		return thumbnail, nil
	}

	// try to find the full-res file.
	thumbnailRoot := strings.TrimSuffix(thumbnail, "_thumbnail.png")

	for _, typ := range allowed_filetypes {
		joinedPath := thumbnailRoot + typ
		if _, err := os.Stat(joinedPath); err == nil {
			// if no error, file exists; return that file.
			return joinedPath, nil
		}
	}

	return thumbnail, nil
	// TODO: create NoFullResolutionErr
	//
	// fmt.Errorf("resolvethumb: no full-res file at path=%s, err=%w", thumbnail)
}

// on my device thumbnails ONLY end with _thumbnail.png or they do not exist.
// this returns the full file path to the highest available resolution of the file.
func GetEagleThumbnailFullRes(cfg *config.NestConfig, itemId string) (string, error) {
	thumbnail, err := GetEagleThumbnail(cfg, itemId)
	if err != nil {
		return thumbnail, fmt.Errorf("getEagleThumbnail: err=%w", err)
	}

	thumbnail, err = resolveThumbnailPath(thumbnail)
	if err != nil {
		return thumbnail, fmt.Errorf("getEagleThumbnail: err=%w", err)
	}

	// TODO: we call os.Stat unnecessarily if we match full-res.
	if _, err = os.Stat(thumbnail); err != nil {
		return thumbnail, fmt.Errorf("getEagleThumbnail: err=%w", err)
	}

	//  TODO: fallback list all files other than metadata.json & _thumbnail.png?
	return thumbnail, nil
}

func GetList(cfg *config.NestConfig) (any, error) {
	// _, err := api.ItemList(cfg.BaseURL(), api.ItemListOptions{Limit: 5})
	_, err := api.ItemList(cfg.BaseURL(), api.ItemListOptions{}.WithDefaults())

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
