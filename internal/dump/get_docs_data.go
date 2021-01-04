/*
Package dump provides functions for saving indexes.
*/
package dump

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"strconv"

	elasticsearch "github.com/elastic/go-elasticsearch/v7"
)

// getDocsData is a function that saves the document of the target index.
func getDocsData(es *elasticsearch.Client, iName string, scrollID string, t int) []map[string]string {
	var (
		buf       bytes.Buffer
		scrollBuf bytes.Buffer
	)

	scrollBuf.WriteString(strconv.Itoa(t))
	scrollBuf.WriteString("m")

	query := map[string]interface{}{
		"scroll":    scrollBuf.String(),
		"scroll_id": scrollID,
	}

	err := json.NewEncoder(&buf).Encode(query)
	if err != nil {
		log.Fatal(err)
	}

	res, err := es.Scroll(
		es.Scroll.WithContext(context.Background()),
		es.Scroll.WithBody(&buf),
		es.Scroll.WithPretty(),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			log.Fatalf("Error parsing the response body: %s", err)
		} else {
			// Print the response status and error information.
			log.Fatalf("[%s] %s: %s",
				res.Status(),
				e["error"].(map[string]interface{})["type"],
				e["error"].(map[string]interface{})["reason"],
			)
		}
	}

	var r map[string]interface{}
	err = json.NewDecoder(res.Body).Decode(&r)
	if err != nil {
		log.Fatal(err)
	}

	var docsData []map[string]string

	for _, hit := range r["hits"].(map[string]interface{})["hits"].([]interface{}) {
		m := map[string]interface{}{
			"index": map[string]interface{}{
				"_index": hit.(map[string]interface{})["_index"],
				"_type":  hit.(map[string]interface{})["_type"],
				"_id":    hit.(map[string]interface{})["_id"],
			},
		}
		index, err := json.Marshal(m)
		if err != nil {
			log.Fatal(err)
		}

		src := hit.(map[string]interface{})["_source"]
		doc, err := json.Marshal(src)
		if err != nil {
			log.Fatal(err)
		}

		docData := map[string]string{
			"index": string(index),
			"doc":   string(doc),
		}

		docsData = append(docsData, docData)
	}

	return docsData
}
