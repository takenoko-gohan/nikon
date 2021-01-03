/*
Package dump provides functions for saving indexes.
*/
package dump

import (
	"fmt"
	"log"

	elasticsearch "github.com/elastic/go-elasticsearch/v7"
)

// SavingIndex is a function that saves the target index to a file.
func SavingIndex(addr string, iName string) {
	cfg := elasticsearch.Config{
		Addresses: []string{
			addr,
		},
	}
	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatal(err)
	}
	scrollID, docsData := getScrollID(es, iName)
	for _, doc := range docsData {
		fmt.Println(doc)
	}
	for i := 0; i < 10; i++ {
		docsData = getDocsData(es, iName, scrollID)
		for _, doc := range docsData {
			fmt.Println(doc)
		}
	}
}
