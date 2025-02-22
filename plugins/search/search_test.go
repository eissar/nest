package search

import (
	"fmt"
	"log"
	"time"

	"github.com/eissar/nest/eagle"
	// "testing"
)

func TestSearch() {
	e, err := eagle.New()
	if err != nil {
		log.Fatalf("%s", err.Error())
	}
	s := New(e)
	defer s.Index.Close()

	//go search.Index(e, s.Index)
	//search.ForceReIndex(e, s.Index)
	ForceReIndexStreaming(e, s.Index)
	return

	start := time.Now()
	s.Query("vallejo")
	fmt.Print("search took: ", time.Since(start))
}
