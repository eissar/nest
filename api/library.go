package api

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"path/filepath"

	"github.com/eissar/nest/api/endpoints"
)

type Folder struct {
	ID               string   `json:"id"`
	Name             string   `json:"name"`
	Description      string   `json:"description"`
	Children         []Folder `json:"children"`
	ModificationTime int64    `json:"modificationTime,omitempty"`
	Tags             []string `json:"tags,omitempty"`
	IconColor        string   `json:"iconColor,omitempty"`
	Password         string   `json:"password,omitempty"`
	PasswordTips     string   `json:"passwordTips,omitempty"`
	CoverID          string   `json:"coverId,omitempty"`
	OrderBy          string   `json:"orderBy,omitempty"`
	SortIncrease     bool     `json:"sortIncrease,omitempty"`
	Icon             string   `json:"icon,omitempty"`
}

type SmartFolder struct {
	ID               string      `json:"id"`
	Icon             string      `json:"icon"`
	Name             string      `json:"name"`
	Description      string      `json:"description"`
	ModificationTime int64       `json:"modificationTime"`
	Conditions       []Condition `json:"conditions"`
	OrderBy          string      `json:"orderBy,omitempty"`
	SortIncrease     bool        `json:"sortIncrease,omitempty"`
}

type Library struct {
	Path string `json:"path"`
	Name string `json:"name"`
}
type LibraryData struct {
	Folders            []Folder      `json:"folders"`
	SmartFolders       []SmartFolder `json:"smartFolders"`
	QuickAccess        []QuickAccess `json:"quickAccess"`
	TagsGroups         []TagsGroup   `json:"tagsGroups"`
	ModificationTime   int64         `json:"modificationTime"`
	ApplicationVersion string        `json:"applicationVersion"`
	Library            Library       `json:"library"`
}

type Condition struct {
	HashKey string `json:"$$hashKey,omitempty"`
	Match   string `json:"match"`
	Rules   []Rule `json:"rules"`
}

type Rule struct {
	HashKey  string `json:"$$hashKey,omitempty"`
	Method   string `json:"method"`
	Property string `json:"property"`
	Value    any    `json:"value"` // Can be []int or string or []string
}

type QuickAccess struct {
	Type string `json:"type"`
	ID   string `json:"id"`
}

type TagsGroup struct {
	ID    string   `json:"id"`
	Name  string   `json:"name"`
	Tags  []string `json:"tags"`
	Color string   `json:"color,omitempty"`
}

// /api/library/switch
func SwitchLibrary(baseURL string, libraryPath string) error {
	ep := endpoints.LibrarySwitch

	uri := baseURL + ep.Path

	libraryPath = filepath.ToSlash(libraryPath)
	body := fmt.Appendf(nil, `{"libraryPath": "%s"}`, libraryPath) //bytes

	req, err := http.NewRequest(ep.Method, uri, bytes.NewReader(body)) // method, url, body
	if err != nil {
		return fmt.Errorf("list: error creating request err=%w", err)
	}

	var a *EagleMessage
	err = InvokeEagleAPIV2(req, &a)
	if err != nil {
		return err
	}

	if a.Status != "success" {
		return fmt.Errorf("response status recieved from eagle was not `success` message=%v", a)
	}

	//fmt.Println(a)

	return nil
}

type LibraryInfoResponse struct {
	Data   LibraryData `json:"data"`
	Status string      `json:"status"`
}

func GetLibraryInfo(baseURL string) (*LibraryInfoResponse, error) {
	var l *LibraryInfoResponse

	ep := endpoints.LibraryInfo
	uri := baseURL + ep.Path

	req, err := http.NewRequest(ep.Method, uri, nil) // method, url, body
	if err != nil {
		return l, fmt.Errorf("list: error creating request err=%w", err)
	}

	err = InvokeEagleAPIV2(req, &l)
	if err != nil {
		return l, fmt.Errorf("error invoking eagle api err=%v", err)
	}

	if l.Status != "success" {
		return l, fmt.Errorf("response status recieved from eagle was not `success` message=%v", l.Status)
	}

	return l, nil
}
func GetIcon(baseURL string) (string, error) {
	var currentLibraryPath string

	ep := endpoints.LibraryIcon

	uri := baseURL + ep.Path

	req, err := http.NewRequest(ep.Method, uri, nil) // method, url, body
	if err != nil {
		return currentLibraryPath, fmt.Errorf("getcurrlibrarypath: error creating request err=%w", err)
	}

	// FIX
	var a *EagleMessage
	err = InvokeEagleAPIV2(req, &a)
	if err != nil {
		return currentLibraryPath, fmt.Errorf("getcurrlibrarypath: error invoking request err=%w", err)
	}

	if v, ok := a.Data.(string); ok {
		currentLibraryPath, err = url.PathUnescape(v)
		if err != nil {
			return currentLibraryPath, fmt.Errorf("getcurrlibrarypath: error parsing path err=%w", err)
		}
	}

	return currentLibraryPath, nil
}

type RecentLibraries struct {
	EagleData
	Libraries []string `json:"data"`
}

type libs []string

// returns []string paths to libraries
// from /api/library/history
func Recent(baseURL string) ([]string, error) {
	ep := endpoints.LibraryHistory

	uri := baseURL + ep.Path

	req, err := http.NewRequest(ep.Method, uri, nil) // method, url, body
	if err != nil {
		return []string{}, fmt.Errorf("recent: error creating request err=%w", err)
	}

	// FIX
	var resp *EagleArray
	err = InvokeEagleAPIV2(req, &resp)
	if err != nil {
		return []string{}, fmt.Errorf("recent: error invoking request err=%w", err)
	}

	return resp.Data, nil
}
