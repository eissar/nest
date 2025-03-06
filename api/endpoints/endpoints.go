package endpoints

// get helpUri of this endpoint.
// TODO: use path.join...
func (e Endpoint) HelpUri() string {
	return "https://api.eagle.cool" + e.Path
}

type Endpoint struct {
	Path   string
	Method string
}

// Application Endpoints
var (
	ApplicationInfo = Endpoint{
		Path:   "/api/application/info",
		Method: "GET",
	}
)

// Folder Endpoints
var (
	FolderCreate = Endpoint{
		Path:   "/api/folder/create",
		Method: "POST",
	}
	FolderRename = Endpoint{
		Path:   "/api/folder/rename",
		Method: "POST",
	}
	FolderUpdate = Endpoint{
		Path:   "/api/folder/update",
		Method: "POST",
	}
	FolderList = Endpoint{
		Path:   "/api/folder/list",
		Method: "GET",
	}
	FolderListRecent = Endpoint{
		Path:   "/api/folder/listRecent",
		Method: "GET",
	}
)

// Item Endpoints
var (
	ItemAddFromURL = Endpoint{
		Path:   "/api/item/addFromURL",
		Method: "POST",
	}
	ItemAddFromURLs = Endpoint{
		Path:   "/api/item/addFromURLs",
		Method: "POST",
	}
	ItemAddFromPath = Endpoint{
		Path:   "/api/item/addFromPath",
		Method: "POST",
	}
	ItemAddFromPaths = Endpoint{
		Path:   "/api/item/addFromPaths",
		Method: "POST",
	}
	ItemAddBookmark = Endpoint{
		Path:   "/api/item/addBookmark",
		Method: "POST",
	}
	ItemInfo = Endpoint{
		Path:   "/api/item/info",
		Method: "GET",
	}
	ItemThumbnail = Endpoint{
		Path:   "/api/item/thumbnail",
		Method: "GET",
	}
	ItemList = Endpoint{
		Path:   "/api/item/list",
		Method: "GET",
	}
	ItemMoveToTrash = Endpoint{
		Path:   "/api/item/moveToTrash",
		Method: "POST",
	}
	ItemRefreshPalette = Endpoint{
		Path:   "/api/item/refreshPalette",
		Method: "POST",
	}
	ItemRefreshThumbnail = Endpoint{
		Path:   "/api/item/refreshThumbnail",
		Method: "POST",
	}
	ItemUpdate = Endpoint{
		Path:   "/api/item/update",
		Method: "POST",
	}
)

// Library Endpoints
var (
	LibraryInfo = Endpoint{
		Path:   "/api/library/info",
		Method: "GET",
	}
	LibraryHistory = Endpoint{
		Path:   "/api/library/history",
		Method: "GET",
	}
	LibrarySwitch = Endpoint{
		Path:   "/api/library/switch",
		Method: "POST",
	}
	LibraryIcon = Endpoint{
		Path:   "/api/library/icon",
		Method: "GET",
	}
)
