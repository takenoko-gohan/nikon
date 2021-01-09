/*
Package indices  provides functions for getting index information.
*/
package indices

import (
	"bytes"
	"log"
	"os"
	"strconv"
	"strings"

	elasticsearch "github.com/elastic/go-elasticsearch/v7"
)

// GetDocCount is a function that gets the number of documents in the target index.
func GetDocCount(es *elasticsearch.Client, iName string) int {
	count := es.Cat.Count
	res, err := count(count.WithIndex(iName), count.WithH("count"))
	if err != nil {
		log.SetOutput(os.Stderr)
		log.Fatal(err)
	}
	defer res.Body.Close()

	buf := new(bytes.Buffer)
	buf.ReadFrom(res.Body)
	result := buf.String()

	result = strings.TrimSpace(result)
	var cnt int
	cnt, _ = strconv.Atoi(result)

	return cnt
}
