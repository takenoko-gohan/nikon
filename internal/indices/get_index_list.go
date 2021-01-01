/*
 */
package indices

import (
	"fmt"
	"log"
	//"encoding/json"

	elasticsearch "github.com/elastic/go-elasticsearch/v7"
)

func GetIndexList(t string) {
	//var b   map[string]interface{}
	fmt.Println(t)
	es, err := elasticsearch.NewDefaultClient()
	if err != nil {
		log.Fatal(err)
	}
	indices := es.Cat.Indices
	res, err := indices(indices.WithV(true))
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	//if err := json.NewDecoder(res.Body).Decode(&b); err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Println(b)
	fmt.Println(res)
}
