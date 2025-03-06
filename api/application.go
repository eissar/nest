package api

import (
	"fmt"
	"net/http"

	"github.com/eissar/nest/api/endpoints"
)

//  TODO:
// - [X] /api/application/info

type ApplicationInfoData struct {
	Version           string `json:"version"`
	PreReleaseVersion string `json:"preReleaseVersion,omitempty"`
	BuildVersion      string `json:"buildVersion"`
	ExecPath          string `json:"execPath"`
	Platform          string `json:"platform"`
}

type ApplicationInfoResponse struct {
	EagleData
	Data ApplicationInfoData `json:"data"`
}

// GET Get detailed information on the Eagle App currently running. In most cases, this could be used to determine whether certain functions are available on the user's device.
// <https://api.eagle.cool/application/info>
func ApplicationInfo(baseUrl string) (ApplicationInfoData, error) {
	ep, ok := endpoints.Application["info"]
	if !ok {
		return ApplicationInfoData{}, fmt.Errorf("could not find endpoint `list` in endpoints.")
	}

	uri := baseUrl + ep.Path

	req, err := http.NewRequest(ep.Method, uri, nil) // method, url, body
	if err != nil {
		return ApplicationInfoData{}, fmt.Errorf("list: error creating request err=%w", err)
	}

	//// add query params
	//query := req.URL.Query()
	// req.URL.RawQuery = query.Encode()
	// fmt.Println("query here:", req.URL.Query().Encode())

	var resp *ApplicationInfoResponse
	err = InvokeEagleAPIV2(req, &resp)
	if err != nil {
		return ApplicationInfoData{}, err
	}
	return resp.Data, nil
}
