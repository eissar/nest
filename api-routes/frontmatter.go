package apiroutes

import (
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	_ "github.com/adrg/frontmatter"
)

func GetNotesCount() string {
	cloudDir := os.Getenv("CLOUD_DIR")
	cloudDir = path.Clean(cloudDir)
	notesPath := path.Join(cloudDir, "Catalog")

	// type WalkDirFunc func(path string, d DirEntry, err error) error
	cnt := 0
	countFn := func(path string, d fs.DirEntry, err error) error {
		if path == notesPath {
			return nil
		}
		if d.IsDir() {
			//fmt.Println(path)
			return filepath.SkipDir
		}

		if strings.HasSuffix(d.Name(), ".md") {
			cnt += 1
		}
		return nil
	}

	err := filepath.WalkDir(notesPath, countFn)
	if err != nil {
		fmt.Println(err)
	}

	return strconv.Itoa(cnt)

}

func GetNotesDetail() []fs.FileInfo {
	cloudDir := os.Getenv("CLOUD_DIR")
	cloudDir = path.Clean(cloudDir)
	notesPath := path.Join(cloudDir, "Catalog")

	notesDetail := []fs.FileInfo{}

	// type WalkDirFunc func(path string, d DirEntry, err error) error
	walkNotes := func(path string, d fs.DirEntry, err error) error {
		// skip root.
		if path == notesPath {
			return nil
		}
		if d.IsDir() {
			//fmt.Println(path)
			return filepath.SkipDir
		}
		if !strings.HasSuffix(d.Name(), ".md") {
			return nil
		}

		i, err := d.Info()
		if err != nil {
			return nil
		}
		notesDetail = append(notesDetail, i)
		return nil // no err
	}

	err := filepath.WalkDir(notesPath, walkNotes)
	if err != nil {
		fmt.Println(err)
	}

	return notesDetail
}

type NotesInfo struct {
	Name         string    `json:"name"`
	ModifiedTime time.Time `json:"modifiedTime"`
}

func GetNotesNamesDates() []NotesInfo {
	a := []NotesInfo{}
	for _, note := range GetNotesDetail() {
		a = append(a, NotesInfo{
			Name:         note.Name(),
			ModifiedTime: note.ModTime(),
		})
	}
	return a
}
