/*
Package dump provides functions for saving indexes.
*/
package dump

import (
	"bytes"
	"encoding/json"
	"log"
	"os"

	elasticsearch "github.com/elastic/go-elasticsearch/v7"
)

/*
deleteScrollID is a function that deletes the scroll_id
passed as an argument from the target Elasticsearch.
*/
func deleteScrollID(es *elasticsearch.Client, scrollID string) {
	var buf bytes.Buffer

	query := map[string]interface{}{
		"scroll_id": scrollID,
	}
	err := json.NewEncoder(&buf).Encode(query)
	if err != nil {
		log.SetOutput(os.Stderr)
		log.Fatal(err)
	}
	res, err := es.ClearScroll(
		es.ClearScroll.WithBody(&buf),
	)
	if err != nil {
		log.SetOutput(os.Stderr)
		log.Fatal(err)
	}
	defer res.Body.Close()

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			log.SetOutput(os.Stderr)
			log.Fatalf("Error parsing the response body: %s", err)
		} else {
			// Print the response status and error information.
			log.SetOutput(os.Stderr)
			log.Fatalf("[%s] %s: %s",
				res.Status(),
				e["error"].(map[string]interface{})["type"],
				e["error"].(map[string]interface{})["reason"],
			)
		}
	}

	log.Println("The scroll id was successfully deleted.")
}
