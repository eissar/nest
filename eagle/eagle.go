package eagle

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sync"
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

func New() (Eagle, error) {
	eag := Eagle{}
	if err := PopulateLibraries(&eag); err != nil {
		return eag, err
	}
	return eag, nil
}
