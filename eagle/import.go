package eagle

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/eissar/nest/config"
	"github.com/eissar/nest/fileUtils"
)

// related to importing libraries into the
// eagle object.

type LibraryImportErr struct {
	ImportError error
	FilePath    string
}

type LibraryNotExistErr struct {
	LibraryImportErr
}

func (e *LibraryImportErr) Error() string {
	return fmt.Sprintf("import: unknown error importing library in libraries.json: library=%s err=%v", e.FilePath, e.ImportError)
}

type LibrariesImportErr struct {
	ImportError error
	FilePath    string
	LenBytes    int // is -1 if error happens
}

func (e *LibrariesImportErr) Error() string {
	const (
		remediationPopulateNestJson     = "You will need to populate it with a JSON array of eagle library paths (folders ending in .library). E.g.,\n\t[\"C:\\Users\\user\\Documents\\Libs\\MyLibrary.library\"]"
		remediationSyntaxJson           = "Make sure the file's contents follow the syntax for a JSON array strictly. You can use an online parser to validate the syntax. The file contents at `path` should look something like this:\n\t[\"C:\\Users\\user\\Documents\\Libs\\MyLibrary.library\"]"
		remediationLibrariesLessThanOne = "Make sure libaries.json has at least one library"
	)
	if errors.Is(e.ImportError, fs.ErrNotExist) {
		return fmt.Sprintf("import: config file does not exist at: path=%s err=%s", e.FilePath, e.ImportError)
	}

	if e.LenBytes == 0 {
		return fmt.Sprintf("import: config file cannot be empty : path=%s fix=%s", e.FilePath, remediationPopulateNestJson)
	}
	_, ok := e.ImportError.(*json.SyntaxError)
	if ok {
		return fmt.Sprintf("import: rerror parsing libraries.json path=%s err=%s fix=%s", e.FilePath, e.ImportError.Error(), remediationSyntaxJson)
	}
	return fmt.Sprintf("import: unknown error importing libraries.json: path=%s err=%v", e.FilePath, e.ImportError)
}

// and returns length of bytes read, and error
// populates v with contents of <configdir>/filename.json
func unmarshalLibraries(p string, v *[]string) (int, error) {
	cfg := filepath.Join(config.GetPath(), p)
	bytes, err := os.ReadFile(cfg)
	if err != nil {
		return -1, fmt.Errorf("populateJson: err=%w", err)
	}
	err = json.Unmarshal(bytes, v)
	if err != nil {
		return len(bytes), fmt.Errorf("populateJson: unmarshalling json err=%w", err)
	}
	return len(bytes), nil
}

func (lib *Library) exists() error {
	_, err := os.Stat(lib.Path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return &LibraryImportErr{
				fmt.Errorf("libraryImportErr: invalid library found in libraries.json err=%w", err),
				lib.Path,
			}
		}
	}
	return nil
}

// fills out Eagle.Libraries
// returns a LibrariesImportErr or nil
func PopulateLibraries(e *Eagle) error {
	var paths []string
	libraries := filepath.Join(config.GetPath(), "libraries.json")
	lenBytes, err := unmarshalLibraries(libraries, &paths)
	if err != nil {
		return &LibrariesImportErr{
			ImportError: err,
			FilePath:    filepath.Join(config.GetConfigPath(), "libraries.json"),
			LenBytes:    -1,
		}
	}
	if len(paths) == 0 {
		return &LibrariesImportErr{
			ImportError: fmt.Errorf("libraries.json is empty or contains no library paths"),
			FilePath:    filepath.Join(config.GetConfigPath(), "libraries.json"),
			LenBytes:    lenBytes,
		}
	}
	// make sure the libraries exist, then add them to the Libraries object.
	for _, path := range paths {
		_, n := filepath.Split(path)
		lib := Library{
			Name: n,
			Path: path,
		}
		if err := lib.exists(); err != nil {
			return err
		} else {
			e.Libraries = append(e.Libraries, lib)
		}
	}
	return nil
}
