/*
Package dump provides functions for saving indexes.
*/
package dump

import (
	//"context"
	"fmt"
	"log"
	"os"
	"sync"

	"golang.org/x/sync/errgroup"

	elasticsearch "github.com/elastic/go-elasticsearch/v7"
	"github.com/takenoko-gohan/nikon/internal/indices"
	"github.com/takenoko-gohan/nikon/internal/processing"
)

// SavingIndex is a function that saves the target index to a file.
func SavingIndex(addr string, iName string, size int, t int, o string) {
	_, err := os.Stat(o)
	if !os.IsNotExist(err) {
		log.Println("The file already exists.")
		os.Exit(0)
	}

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

	var eg errgroup.Group

	chRes := make(chan map[string]interface{}, 10) 
	chDoc := make(chan []map[string]string, 10)

	var scrollID string

	eg.Go(func() error {
		scrollID, err = getScrollID(es, iName, size, t, chRes)
		if err != nil {
			return err
		}

		return nil
	})

	var mu1 sync.Mutex

	for i := 0; i < cnt; i++ {
		eg.Go(func() error {
			mu1.Lock()
			defer mu1.Unlock()
			return getScrollRes(es, iName, scrollID, t, chRes)
		})
	}

	eg.Go(func() error {
		return processing.GetDocsData(chRes, chDoc)
	})

	var mu2 sync.Mutex

	eg.Go(func() error {
		mu2.Lock()
		defer mu2.Unlock()
		return saveDocToFile(o, chDoc)
	})

	if err := eg.Wait(); err != nil{
		log.Fatal(err)
	}

	deleteScrollID(es, scrollID)
}
