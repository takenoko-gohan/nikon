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

const scrollPrefix = "m"
const getResLogPrefix1 = "got "
const getResLogPrefix2 = " documents from Elasticsearch."

// getScrollRes is function to get the response from the Scroll API
func getScrollRes(es *elasticsearch.Client, iName string, scrollID string, t int, out chan<- map[string]interface{}, done chan<- struct{}) error {
	var (
		buf bytes.Buffer
	)

	scroll := processing.StringConcat([]interface{}{
		t,
		scrollPrefix,
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
	if cnt.Len() == 0 {
		done <- struct{}{}
		return nil
	}

	msg := processing.StringConcat([]interface{}{
		getResLogPrefix1,
		cnt.Len(),
		getResLogPrefix2,
	})
	log.Println(msg)

	out <- r

	return nil
}
