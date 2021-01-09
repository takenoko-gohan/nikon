/*
Package dump provides functions for saving indexes.
*/
package dump

import (
	//"context"
	"context"
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
	log.SetOutput(os.Stdout)

	_, err := os.Stat(o)
	if !os.IsNotExist(err) {
		log.SetOutput(os.Stderr)
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
		log.SetOutput(os.Stderr)
		log.Fatal(err)
	}

	cnt := indices.GetDocCount(es, iName)
	if cnt != 0 {
		cnt = cnt/size + 1
	} else {
		log.SetOutput(os.Stderr)
		log.Println("The document does not exist in the target index.")
		os.Exit(0)
	}

	eg, ctx := errgroup.WithContext(context.Background())

	chRes := make(chan map[string]interface{}, 10)
	chResDone := make(chan struct{})
	chDoc := make(chan []map[string]string, 10)

	var scrollID string

	scrollID, err = getScrollID(es, iName, size, t, chRes)
	if err != nil {
		log.SetOutput(os.Stderr)
		log.Fatal(err)
	}

	var mu1 sync.Mutex

	for i := 0; i < cnt; i++ {
		eg.Go(func() error {
			select {
			case <-ctx.Done():
				return nil
			default:
				mu1.Lock()
				defer mu1.Unlock()
				err := getScrollRes(es, iName, scrollID, t, chRes, chResDone)
				if err != nil {
					return err
				}

				return nil
			}
		})
	}

	eg.Go(func() error {
		select {
		case <-ctx.Done():
			return nil
		case <-chResDone:
			close(chRes)
			return nil
		}
	})

	eg.Go(func() error {
		select {
		case <-ctx.Done():
			return nil
		default:
			defer close(chDoc)
			return processing.GetDocsData(chRes, chDoc)
		}
	})

	var mu2 sync.Mutex

	eg.Go(func() error {
		select {
		case <-ctx.Done():
			return nil
		default:
			mu2.Lock()
			defer mu2.Unlock()
			return saveDocToFile(o, chDoc)
		}
	})

	if err := eg.Wait(); err != nil {
		log.SetOutput(os.Stderr)
		log.Fatal(err)
	}

	deleteScrollID(es, scrollID)

	log.Println("The index was saved successfully.")
}
