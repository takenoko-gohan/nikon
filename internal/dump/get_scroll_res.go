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

	elasticsearch "github.com/elastic/go-elasticsearch/v7"
	"github.com/takenoko-gohan/nikon/internal/processing"
)

// getScrollRes is function to get the response from the Scroll API
func getScrollRes(es *elasticsearch.Client, iName string, scrollID string, t int, out chan<- map[string]interface{}) error {
	var (
		buf bytes.Buffer
	)

	scroll := processing.StringConcat([]interface{}{
		t,
		"m",
	})

	query := map[string]interface{}{
		"scroll":    scroll,
		"scroll_id": scrollID,
	}

	err := json.NewEncoder(&buf).Encode(query)
	if err != nil {
		return err
	}

	res, err := es.Scroll(
		es.Scroll.WithContext(context.Background()),
		es.Scroll.WithBody(&buf),
		es.Scroll.WithPretty(),
	)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		var e map[string]interface{}
		err := json.NewDecoder(res.Body).Decode(&e)
		if err != nil {
			return err
		}
	}

	var r map[string]interface{}
	err = json.NewDecoder(res.Body).Decode(&r)
	if err != nil {
		return err
	}

	cnt := reflect.ValueOf(r["hits"].(map[string]interface{})["hits"])
	msg := processing.StringConcat([]interface{}{
		"got ",
		cnt.Len(),
		" documents from Elasticsearch.",
	})
	log.Println(msg)

	out <- r

	return nil
}
