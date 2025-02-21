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
	"info": { // returns status, data
		Path:   "/api/application/info", // https://api.eagle.cool/application/info
		Method: "GET",
	},
}

var Folder = map[string]Endpoint{
	"create": { // returns status, data
		Path:   "/api/folder/create", // https://api.eagle.cool/folder/create
		Method: "POST",
	},
	"rename": { // returns status, data
		Path:   "/api/folder/rename", // https://api.eagle.cool/folder/rename
		Method: "POST",
	},
	"update": { // returns status, data
		Path:   "/api/folder/update", // https://api.eagle.cool/folder/update
		Method: "POST",
	},
	"list": { // returns status, data
		Path:   "/api/folder/list", // https://api.eagle.cool/folder/list
		Method: "GET",
	},
	"listRecent": { // returns status, data
		Path:   "/api/folder/listRecent", // https://api.eagle.cool/folder/listRecent
		Method: "GET",
	},
}

//eagle.cool/api/item/ endpoints
var Item = map[string]Endpoint{
	"addFromURL": { // returns status
		Path:   "/api/item/addFromURL", // https://api.eagle.cool/item/add-from-uRL
		Method: "POST",
	},
	"addFromURLs": { // returns status
		Path:   "/api/item/addFromURLs", // https://api.eagle.cool/item/add-From-URLs
		Method: "POST",
	},
	"addFromPath": { // returns status
		Path:   "/api/item/addFromPath", // https://api.eagle.cool/item/addFromPath
		Method: "POST",
	},
	"addFromPaths": { // returns status
		Path:   "/api/item/addFromPaths", // https://api.eagle.cool/item/addFromPaths
		Method: "POST",
	},
	"addBookmark": { // returns status
		Path:   "/api/item/addBookmark", // https://api.eagle.cool/item/addBookmark
		Method: "POST",
	},
	"info": { // returns status, data
		Path:   "/api/item/info", // https://api.eagle.cool/item/info
		Method: "GET",
	},
	"thumbnail": { // returns status, data
		Path:   "/api/item/thumbnail", // https://api.eagle.cool/item/thumbnail
		Method: "GET",
	},
	"list": { // returns status, data
		Path:   "/api/item/list", // https://api.eagle.cool/item/list
		Method: "GET",
	},
	"moveToTrash": { // returns status
		Path:   "/api/item/moveToTrash", // https://api.eagle.cool/item/moveToTrash
		Method: "POST",
	},
	"refreshPalette": { // returns status
		Path:   "/api/item/refreshPalette", // https://api.eagle.cool/item/refreshPalette
		Method: "POST",
	},
	"refreshThumbnail": { // returns status
		Path:   "/api/item/refreshThumbnail", // https://api.eagle.cool/item/refreshThumbnail
		Method: "POST",
	},
	"update": { // returns status, data
		Path:   "/api/item/update", // https://api.eagle.cool/item/update
		Method: "POST",
	},
}

var Library = map[string]Endpoint{
	"info": { // returns status, data
		Path:   "/api/library/info", // https://api.eagle.cool/library/info
		Method: "GET",
	},
	"history": { // returns status, data
		Path:   "/api/library/history", // https://api.eagle.cool/library/history
		Method: "GET",
	},
	"switch": { // returns status ( does not wait )
		Path:   "/api/library/switch", // https://api.eagle.cool/library/switch
		Method: "POST",
	},
	"icon": { // returns imageBytes in request body.
		Path:   "/api/library/icon", // https://api.eagle.cool/library/icon
		Method: "GET",
	},
}
