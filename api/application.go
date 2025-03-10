package api

import (
	"fmt"

	"github.com/eissar/nest/api/endpoints"
)

//	TODO:
//
// - [X] /api/application/info

type ApplicationInfoData struct {
	Version           string `json:"version"`
	PreReleaseVersion string `json:"preReleaseVersion,omitempty"`
	BuildVersion      string `json:"buildVersion"`
	ExecPath          string `json:"execPath"`
	Platform          string `json:"platform"`
}

// GET Get detailed information on the Eagle App currently running. In most cases, this could be used to determine whether certain functions are available on the user's device.
// <https://api.eagle.cool/application/info>
func ApplicationInfo(baseUrl string) (ApplicationInfoData, error) {
	ep := endpoints.ApplicationInfo
	uri := baseUrl + ep.Path

	var resp struct {
		EagleResponse                     // Response string `json:"response"`
		Data          ApplicationInfoData `json:"data"`
	}
	err := Request(ep.Method, uri, nil, nil, &resp)
	if err != nil {
		return resp.Data, fmt.Errorf("ApplicationInfo: err=%w", err)
	}
	if resp.Status != "success" {
		return resp.Data, fmt.Errorf("ApplicationInfo: err=%w", ErrStatusErr)
	}

	return resp.Data, nil
}
