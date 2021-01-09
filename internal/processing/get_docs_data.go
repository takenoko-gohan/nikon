/*
Package processing provides functions to
process the values passed in.
*/
package processing

import (
	"encoding/json"
)

// GetDocsData is a function that retrieves information used
// by the Bulk API from Elasticsearch responses.
func GetDocsData(in <-chan map[string]interface{}, out chan<- []map[string]string) error {
	for r := range in {
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
				return err
			}

			src := hit.(map[string]interface{})["_source"]
			doc, err := json.Marshal(src)
			if err != nil {
				return err
			}

			docData := map[string]string{
				"index": string(index),
				"doc":   string(doc),
			}

			docsData = append(docsData, docData)
		}

		out <- docsData
	}

	return nil
}
