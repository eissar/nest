package search

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/blevesearch/bleve/v2"
	"github.com/eissar/nest/config"
	"github.com/eissar/nest/eagle"
	"github.com/eissar/nest/fileUtils"
)

type Search struct {
	eagle.Eagle
	Index bleve.Index
}

// create or open a new index
// mapping := bleve.NewIndexMapping()
// remember to call Search.Index.Close.
// may call log.Fatalf
func New(e eagle.Eagle) *Search {
	newIndexFlag := false
	blevePath := filepath.Join(config.GetPath(), "example.bleve")
	_, err := os.Stat(blevePath)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			newIndexFlag = true
		} else {
			log.Fatalf("search.initialize: error=%v", err.Error())
		}
	}

	// use existing ...
	if newIndexFlag == false {
		i, err := bleve.Open(blevePath)
		if err != nil {
			panic(err)
		}
		return &Search{
			Eagle: e,
			Index: i,
		}
	}

	// make a new one...
	mapping := bleve.NewIndexMapping()
	eagleImageMapping := bleve.NewDocumentMapping()

	lastModifiedFieldMapping := bleve.NewNumericFieldMapping()
	//lastModifiedFieldMapping.Type = "int"

	eagleImageMapping.AddFieldMappingsAt("lastModified", lastModifiedFieldMapping)

	mapping.AddDocumentMapping("Image", eagleImageMapping)

	i, err := bleve.New(blevePath, mapping)
	if err != nil {
		panic(err)
	}
	return &Search{
		Eagle: e,
		Index: i,
	}
}

// create or open a new index
// mapping := bleve.NewIndexMapping()
func Initialize(e eagle.Eagle) bleve.Index {
	newIndexFlag := false
	blevePath := filepath.Join(config.GetPath(), "example.bleve")
	_, err := os.Stat(blevePath)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			newIndexFlag = true
		} else {
			log.Fatalf("search.initialize: error=%v", err.Error())
		}
	}

	if newIndexFlag == false { // use existing ...
		i, err := bleve.Open(blevePath)
		if err != nil {
			panic(err)
		}
		return i
	}
	// make a new one...

	mapping := bleve.NewIndexMapping()
	eagleImageMapping := bleve.NewDocumentMapping()

	lastModifiedFieldMapping := bleve.NewNumericFieldMapping()
	//lastModifiedFieldMapping.Type = "int"

	eagleImageMapping.AddFieldMappingsAt("lastModified", lastModifiedFieldMapping)

	mapping.AddDocumentMapping("Image", eagleImageMapping)

	i, err := bleve.New(blevePath, mapping)
	if err != nil {
		panic(err)
	}
	return i
}

// search for some text
func (s *Search) Query(q string) {
	if q == "" {
		log.Fatal("q string cannot be empty in call to Query.")
	}
	// todo add a resolver to map
	// folderIds to folder names.

	//query := bleve.NewMatchQuery("vallejo")
	query := bleve.NewQueryStringQuery("*vallejo*")
	//query := bleve.NewFuzzyQuery(q)
	//query.SetFuzziness(2)
	search := bleve.NewSearchRequest(query)
	search.Fields = []string{"*"}
	searchResults, err := s.Index.Search(search)
	if err != nil {
		panic(err)
	}

	for _, hit := range searchResults.Hits {
		fmt.Printf("searchResults: %v\n", hit.Fields)
		fmt.Printf("searchResults: %v\n", hit.Score)
	}
}

func Index(e eagle.Eagle, i bleve.Index) {
	start := time.Now()

	lib := e.Libraries[0]
	imgDirs := lib.WalkImageDirs()

	imagesBatch := i.NewBatch()
	fmt.Printf("iterating... \n")
	for _, dir := range imgDirs {
		var item eagle.Image
		err := eagle.ParseItemMetadata(dir, &item)
		if err != nil {
			continue
		}

		//bleve.NewQueryStringQuery(")
		// query func () *query.ConjunctionQuery {
		// 	q1 := bleve.NewDocIDQuery([]string{item.Id})
		// 	bleve.NewConjunctionQuery()

		// }()
		q := bleve.NewDocIDQuery([]string{item.Id})
		search := bleve.NewSearchRequest(q)
		search.Fields = []string{"lastModified", "id"}
		searchResults, err := i.Search(search)
		if err != nil {
			panic(err)
		}

		type result struct {
			Id           string `json:"id"`
			LastModified int    `json:"lastModified"`
		}

		// return true to add item.
		addFlag := func() bool {
			// localStart := time.Now()
			// defer fmt.Printf("%v", time.Since(localStart))
			for _, hit := range searchResults.Hits {
				mod, ok := hit.Fields["lastModified"].(float64)
				if !ok {
					fmt.Printf("failed type assertion for value: %v to float64\n", hit.Fields["lastModified"])
					return true
				}

				// item.LastModified is the most up-to-date.
				// (from metadata.json)

				// only re-index items whose mtime are less than < item.LastModified
				if int(mod) < item.LastModified {
					return true
				}
			}
			// no hits; return true
			return true
		}()
		//index.Batch()

		if addFlag {
			imagesBatch.Index(item.Id, item)
		}
	}
	fmt.Printf("finished loading inserts in %v\n", time.Since(start))

	if err := i.Batch(imagesBatch); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("indexing completed in %s\n", time.Since(start))
}
func ForceReIndex(e eagle.Eagle, i bleve.Index) {
	start := time.Now()

	lib := e.Libraries[0]
	imgDirs := lib.WalkImageDirs()
	fmt.Printf("finished loading image dirs in %v\n", time.Since(start))

	imagesBatch := i.NewBatch()
	for _, dir := range imgDirs {
		var item eagle.Image
		err := eagle.ParseItemMetadata(dir, &item)
		if err != nil {
			// log
			continue
		}
		imagesBatch.Index(item.Id, item)
	}
	fmt.Printf("finished loading inserts in %v\n", time.Since(start))

	if err := i.Batch(imagesBatch); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("indexing completed in %s\n", time.Since(start))
}
func ForceReIndexStreaming(e eagle.Eagle, i bleve.Index) {
	start := time.Now()

	lib := e.Libraries[0]
	imgDirs := lib.WalkImageDirs()
	fmt.Printf("finished loading image dirs in %v\n", time.Since(start))

	readers := make([]io.Reader, 0, len(imgDirs)) // Pre-allocate
	for _, dir := range imgDirs {
		handle, err := os.Open(dir)
		if err != nil {
			continue
		}
		defer handle.Close()
		readers = append(readers, handle)
	}
	fmt.Printf("hndls: %v\n", len(readers))

	multiReader := io.MultiReader(readers...)
	dcdr := json.NewDecoder(multiReader)
	imagesBatch := i.NewBatch()

	for {
		var item eagle.Image
		err := dcdr.Decode(&item)
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println("Error decoding JSON:", err)
			return
		}

		imagesBatch.Index(item.Id, &item)
		//fmt.Println(item.Id)
	}
	fmt.Printf("finished loading inserts in %v\n", time.Since(start))

	if err := i.Batch(imagesBatch); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("indexing completed in %s\n", time.Since(start))
}
