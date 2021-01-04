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
func SavingIndex(addr string, iName string, size int) {
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

	scrollID, docsData := getScrollID(es, iName, size)
	for _, doc := range docsData {
		fmt.Println(doc["index"])
		fmt.Println(doc["doc"])
	}
	for i := 0; i < cnt; i++ {
		docsData = getDocsData(es, iName, scrollID)
		for _, doc := range docsData {
			fmt.Println(doc["index"])
			fmt.Println(doc["doc"])
		}
	}
	deleteScrollID(es, scrollID)
}
