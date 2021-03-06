/*
Package indices  provides functions for getting index information.
*/
package indices

import (
	"bytes"
	"fmt"
	"log"
	"os"

	elasticsearch "github.com/elastic/go-elasticsearch/v7"
)

// GetIndexList is a function to get the index list of Elasticsearch specified by the argument.
func GetIndexList(addr string) {
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
	indices := es.Cat.Indices
	res, err := indices(indices.WithH("index"))
	if err != nil {
		log.SetOutput(os.Stderr)
		log.Fatal(err)
	}
	defer res.Body.Close()

	buf := new(bytes.Buffer)
	buf.ReadFrom(res.Body)
	result := buf.String()

	fmt.Println(result)
}
