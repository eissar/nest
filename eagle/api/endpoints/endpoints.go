package endpoints

// get helpUri of this endpoint.
func (e Endpoint) HelpUri() string {
	return "https://api.eagle.cool" + e.Path
}

type Endpoint struct {
	Path   string
	Method string
}

// maps defining endpoints

var Application = map[string]Endpoint{
	"info": {
		Path:   "/api/application/info",
		Method: "POST",
	},
}

var Folder = map[string]Endpoint{
	"create": {
		Path:   "/api/folder/create",
		Method: "POST",
	},
	"rename": {
		Path:   "/api/folder/rename",
		Method: "POST",
	},
	"update": {
		Path:   "/api/folder/update",
		Method: "POST",
	},
	"list": {
		Path:   "/api/folder/list",
		Method: "GET",
	},
	"listRecent": {
		Path:   "/api/folder/listRecent",
		Method: "GET",
	},
}

var ItemEndpoints = map[string]Endpoint{
	"addFromURL": {
		Path:   "/api/item/addFromURL",
		Method: "POST",
	},
	"addFromURLs": {
		Path:   "/api/item/addFromURLs",
		Method: "POST",
	},
	"addFromPath": {
		Path:   "/api/item/addFromPath",
		Method: "POST",
	},
	"addFromPaths": {
		Path:   "/api/item/addFromPaths",
		Method: "POST",
	},
	"addBookmark": {
		Path:   "/api/item/addBookmark",
		Method: "POST",
	},
	"info": {
		Path:   "/api/item/info",
		Method: "GET",
	},
	"thumbnail": {
		Path:   "/api/item/thumbnail",
		Method: "GET",
	},
	"list": {
		Path:   "/api/item/list",
		Method: "GET",
	},
	"moveToTrash": {
		Path:   "/api/item/moveToTrash",
		Method: "POST",
	},
	"refreshPalette": {
		Path:   "/api/item/refreshPalette",
		Method: "POST",
	},
	"refreshThumbnail": {
		Path:   "/api/item/refreshThumbnail",
		Method: "POST",
	},
	"update": {
		Path:   "/api/item/update",
		Method: "POST",
	},
}

var Library = map[string]Endpoint{
	"info": {
		Path:   "/api/library/info",
		Method: "GET",
	},
	"history": {
		Path:   "/api/library/history",
		Method: "GET",
	},
	"switch": {
		Path:   "/api/library/switch",
		Method: "POST",
	},
	"icon": {
		Path:   "/api/library/icon",
		Method: "GET",
	},
}
