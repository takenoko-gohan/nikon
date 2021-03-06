/*
Package dump provides functions for saving indexes.
*/
package dump

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"reflect"
	"time"

	elasticsearch "github.com/elastic/go-elasticsearch/v7"
	"github.com/takenoko-gohan/nikon/internal/processing"
)

const getScrollIDLogPrefix1 = "got "
const getScrollIDLogPrefix2 = " documents from Elasticsearch."

// getScrollID is a function that gets scroll_id and the first document.
func getScrollID(es *elasticsearch.Client, iName string, size int, t int, out chan<- map[string]interface{}) (string, error) {
	var buf bytes.Buffer

	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match_all": map[string]interface{}{},
		},
		"size": size,
	}

	err := json.NewEncoder(&buf).Encode(query)
	if err != nil {
		return "", err
	}

	scrollT := t * 60000000000

	res, err := es.Search(
		es.Search.WithContext(context.Background()),
		es.Search.WithIndex(iName),
		es.Search.WithBody(&buf),
		es.Search.WithScroll(time.Duration(scrollT)),
		es.Search.WithPretty(),
	)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			return "", err
		}
	}

	var r map[string]interface{}
	err = json.NewDecoder(res.Body).Decode(&r)
	if err != nil {
		return "", err
	}

	cnt := reflect.ValueOf(r["hits"].(map[string]interface{})["hits"])
	msg := processing.StringConcat([]interface{}{
		getScrollIDLogPrefix1,
		cnt.Len(),
		getScrollIDLogPrefix2,
	})
	log.Println(msg)

	out <- r

	scrollID := r["_scroll_id"].(string)

	return scrollID, nil
}
