package eagle

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
	"web-dashboard/config"
	"web-dashboard/fileUtils"
)

type Library struct {
	Name  string
	Path  string
	Mutex sync.Mutex
}
type Eagle struct {
	//Db        *sql.DB
	Libraries []*Library
}

type Image struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Size int    `json:"size"`
	Ext  string `json:"ext"`
}

// Ptr is a type constraint for pointers to any type.
type Ptr[T any] interface{ *T }

// mustUnmarshalFromFile unmarshals data from a file into the provided pointer.
// it panics if path does not exist.
func MustParseJson[T any, P Ptr[T]](pth string, v P) error {
	bytes, err := os.ReadFile(pth)
	if err != nil {
		log.Fatalf("must unmarshal: unknown import error reading from %s err: %s", pth, err)
	}

	err = json.Unmarshal(bytes, v)
	if err != nil {
		log.Fatalf("must unmarshal: import error while unmarshalling %s err: %s", pth, err)
	}
	return nil
}

func (lib *Library) EnumJson() {
	lib.Mutex.Lock()
	defer lib.Mutex.Unlock()

	start := time.Now()
	pth := filepath.Join(lib.Path, "images")

	totalCnt := 0
	pathsArr := []string{}
	filepath.WalkDir(pth, func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}
		// d.Name()
		if d.Name() == "metadata.json" {
			totalCnt += 1
			//if totalCnt < 10 {
			pathsArr = append(pathsArr, path)
			//}
			return nil
		}
		return nil
	})
	//fmt.Println(totalCnt)

	maxLen := len(pathsArr)
	for i, path := range pathsArr {
		fmt.Printf("Progress: %d \r", (i*100)/maxLen)
		img := &Image{}
		MustParseJson(path, img)
		//fmt.Println(img)
	}
	fmt.Println("finished parsing json in: ", time.Since(start))
}

func (lib *Library) FirstFiveImages() []*Image {
	lib.Mutex.Lock()
	defer lib.Mutex.Unlock()

	pth := filepath.Join(lib.Path, "images")

	totalCnt := 0
	pathsArr := []string{}
	filepath.WalkDir(pth, func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}
		// d.Name()
		if d.Name() == "metadata.json" {
			totalCnt += 1
			//if totalCnt < 10 {
			pathsArr = append(pathsArr, path)
			//}
			return nil
		}
		return nil
	})
	//fmt.Println(totalCnt)

	var images []*Image
	//maxLen := len(pathsArr)
	for _, path := range pathsArr[0:6] {
		img := &Image{}
		MustParseJson(path, img)
		//fmt.Println(img)
		images = append(images, img)
	}
	return images
}

func New() *Eagle {
	eag := &Eagle{}
	return eag
}

type LibraryImportErr error

type LibrariesImportErr struct {
	ImportError error
	FilePath    string
	LenBytes    int // is -1 if error happens
}

func (e *LibrariesImportErr) Error() string {
	remediationPopulateNestJson := "You will need to populate it " +
		"with a json array of eagle library paths (folders ending in .library) E.g.,\n\t" +
		`["C:\Users\user\Documents\Libs\MyLibrary.library"]`
	remediationSyntaxJson := "Make sure the file's contents follow the " +
		"syntax for a json array strictly. you can use an online parser to validate the syntax." +
		"the file contents at `path` should look something like this:\n\t" +
		`["C:\Users\user\Documents\Libs\MyLibrary.library"]`

	if errors.Is(e.ImportError, fs.ErrNotExist) {
		msg := "import: config file does not exist at: path=%s this should not have happened. report please! err=%s"
		return fmt.Sprintf(msg, e.FilePath, e.ImportError)
	}

	if e.LenBytes == 0 {
		msg := "import: config file cannot be empty : path=%s fix=%s"
		return fmt.Sprintf(msg, e.FilePath, remediationPopulateNestJson)
	}
	_, ok := e.ImportError.(*json.SyntaxError)
	if ok {
		msg := "import: error parsing config : path=%s err=%s fix=%s"
		return fmt.Sprintf(msg, e.FilePath, e.ImportError.Error(), remediationSyntaxJson)
	}
	return fmt.Sprintf("import: unknown error importing config.json: path=%s err=%s", e.FilePath, e.ImportError.Error())
}

// fills out Eagle.Libraries
// returns a LibrariesImportErr or nil
func PopulateLibraries(e *Eagle) error {
	var paths []string
	lenBytes, err := config.PopulateJson("libraries.json", &paths)
	if err != nil {
		return &LibrariesImportErr{
			ImportError: err,
			FilePath:    filepath.Join(config.GetConfigPath(), "libraries.json"),
			LenBytes:    -1,
		}
	}
	if len(paths) == 0 {
		return &LibrariesImportErr{
			ImportError: fmt.Errorf("import error after parsing config.json - json " +
				"empty or malformed? no library paths were found."),
			FilePath: filepath.Join(config.GetConfigPath(), "libraries.json"),
			LenBytes: lenBytes,
		}
	}
	// make sure the paths exists, then add them to the Libraries object.
	for _, lib := range paths {
		lib = filepath.Clean(lib)

		err := fileUtils.TestPath(lib)
		if err != nil {
			if errors.Is(err, fs.ErrNotExist) {
				return LibraryImportErr(fmt.Errorf("libraryImportErr: invalid library found "+
					"in libraries.json err=%w", err))
			}
		} else {
			return LibraryImportErr(fmt.Errorf("libraryImportErr: unknown error "+
				"err=%w", err))
		}

		_, n := filepath.Split(lib)
		e.Libraries = append(e.Libraries, &Library{
			Name: n,
			Path: lib,
		})
	}
	return nil
}
