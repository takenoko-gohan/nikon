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
	"github.com/takenoko-gohan/nikon/internal/processing"
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
	msg := processing.StringConcat([]interface{}{
		"got ",
		size,
		" documents from Elasticsearch.(offset: 0)",
	})
	log.Println(msg)
	saveDocToFile(o, docsData, true)
	msg = processing.StringConcat([]interface{}{
		"saved ",
		size,
		" documents to a file.(offset: 0)",
	})
	log.Println(msg)
	for i := 0; i < cnt; i++ {
		offset := (i + 1) * size
		docsData = getDocsData(es, iName, scrollID, t)
		processing.StringConcat([]interface{}{
			"got ",
			size,
			" documents from Elasticsearch.(offset: ",
			offset,
			")",
		})
		log.Println(msg)
		saveDocToFile(o, docsData, false)
		msg = processing.StringConcat([]interface{}{
			"saved ",
			size,
			" documents to a file.(offset: ",
			offset,
			")",
		})
		log.Println(msg)
	}
	deleteScrollID(es, scrollID)
}
