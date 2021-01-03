/*
Package indices  provides functions for getting index information.
*/
package indices

import (
	"fmt"
	"log"
	"bytes"

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
		log.Fatal(err)
	}
	indices := es.Cat.Indices
	res, err := indices(indices.WithV(true))
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	buf := new(bytes.Buffer)
	buf.ReadFrom(res.Body)
	resutl := buf.String()

	fmt.Println(resutl)
}
