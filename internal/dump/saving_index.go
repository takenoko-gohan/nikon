/*
Package dump provides functions for saving indexes.
*/
package dump

import (
	"fmt"
	"log"
	"os"

	elasticsearch "github.com/elastic/go-elasticsearch/v7"
	"github.com/takenoko-gohan/nikon/internal/indices"
)

// SavingIndex is a function that saves the target index to a file.
func SavingIndex(addr string, iName string, size int, t int, o string) {
	cfg := elasticsearch.Config{
		Addresses: []string{
			addr,
		},
	}
	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatal(err)
	}

	cnt := indices.GetDocCount(es, iName)
	if cnt != 0 {
		cnt = cnt / size
	} else {
		fmt.Println("The document does not exist in the target index.")
		os.Exit(0)
	}

	scrollID, docsData := getScrollID(es, iName, size, t)
	log.Println("got ", size, " documents from Elasticsearch.(offset: 0)")
	saveDocToFile(o, docsData, true)
	log.Println("saved ", size, " documents to a file.(offset: 0)")
	for i := 0; i < cnt; i++ {
		offset := (i + 1) * size
		docsData = getDocsData(es, iName, scrollID, t)
		log.Println("got ", size, " documents from Elasticsearch.(offset: ", offset, ")")
		saveDocToFile(o, docsData, false)
		log.Println("saved ", size, " documents to a file.(offset: ", offset, ")")
	}
	deleteScrollID(es, scrollID)
}
