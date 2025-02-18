package eagle

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Library struct {
	Name string
	Path string
}
type Eagle struct {
	Libraries []Library
	//Db        *sql.DB
}

type Image struct {
	Id           string   `json:"id"`
	Name         string   `json:"name"`
	Size         int      `json:"size"`
	Ext          string   `json:"ext"`
	LastModified int      `json:"lastModified"` //eagleModTime
	Tags         []string `json:"tags"`
	Url          string   `json:"url"`
	Annotation   string   `json:"annotation"`
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

func New() (Eagle, error) {
	eag := Eagle{}
	if err := PopulateLibraries(&eag); err != nil {
		return eag, err
	}
	return eag, nil
}

// get metadata.json
func (lib *Library) WalkImageDirs() (items []string) {
	imagesDir := filepath.Join(lib.Path, "images")
	imagesDirSeps := strings.Count(imagesDir, string(filepath.Separator))
	filepath.WalkDir(imagesDir, func(path string, d fs.DirEntry, err error) error {
		depth := strings.Count(path, string(filepath.Separator)) - imagesDirSeps
		if depth == 0 {
			// root
			return nil
		} else if depth == 1 {
			items = append(items, filepath.Join(path, "metadata.json"))
			return fs.SkipDir
		}
		return fs.SkipDir
	})
	return items
}

func ParseItemMetadata[T any, P Ptr[T]](pth string, v P) error {
	bytes, err := os.ReadFile(pth)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("ParseItemMetadata: %w", err)
		} else {
			return fmt.Errorf("ParseItemMetadata: %w", err)
		}
	}

	err = json.Unmarshal(bytes, v)
	if err != nil {
		return fmt.Errorf("unmarshalling json ParseItemMetadata: err=%s", err)
	}
	return nil
}
