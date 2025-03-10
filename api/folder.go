package api

// - [ ] /api/folder/create
// - [ ] /api/folder/rename
// - [ ] /api/folder/update
// - [ ] /api/folder/list
// - [ ] /api/folder/listRecent

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/eissar/nest/api/endpoints"
)

// todo rename
type FolderCreateResponse struct {
	ID               string `json:"id"`
	Name             string `json:"name"`
	ModificationTime int    `json:"modificationTime"`
	//Images           ...          `json:"images"`
	//Folders          ...          `json:"folders"`
	//ImagesMappings   ... `json:"imagesMappings"`
	//Tags             ...          `json:"tags"`
	//Children         ...          `json:"children"`
	//IsExpand         bool   `json:"isExpand"`
}

/* { "status": "success",
    "data": {
        "id": "KBJJSMMVF9WYL",
        "name": "The Folder Name",
        "images": [],
        "folders": [],
        "modificationTime": 1592409993367,
        "imagesMappings": {},
        "tags": [],
        "children": [],
        "isExpand": true
    }
}
*/

func FolderCreate(baseUrl string, folderName string) (FolderCreateResponse, error) {
	ep := endpoints.FolderCreate
	uri := baseUrl + ep.Path

	var resp struct {
		EagleResponse
		Data FolderCreateResponse `json:"data"`
	}

	requestBody := struct {
		FolderName string `json:"folderName"`
	}{folderName}

	body, err := json.Marshal(requestBody)
	if err != nil {
		return resp.Data, fmt.Errorf("addfromurl: error converting request into json body err=%w", err)
	}

	err = Request(ep.Method, uri, bytes.NewReader(body), nil, &resp)
	if err != nil {
		return resp.Data, fmt.Errorf("addFromUrl: err=%w", err)
	}

	if resp.Status != "success" {
		return resp.Data, fmt.Errorf("addFromUrl: err=%w", ErrStatusErr)
	}

	return FolderCreateResponse{}, nil
}
func FolderRename()     {}
func FolderList()       {}
func FolderListRecent() {}
