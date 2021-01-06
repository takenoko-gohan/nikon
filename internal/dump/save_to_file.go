/*
Package dump provides functions for saving indexes.
*/
package dump

import (
	"log"
	"os"

	"github.com/takenoko-gohan/nikon/internal/processing"
)

// saveDocToFile is a function that saves the passed document data a file.
func saveDocToFile(o string, docsData []map[string]string, isFirst bool) {
	if isFirst {
		f, err := os.OpenFile(o, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0666)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		for _, doc := range docsData {
			str := processing.StringConcat([]interface{}{
				doc["index"],
				"\n",
				doc["doc"],
				"\n",
			})
			f.WriteString(str)
		}
	} else {
		f, err := os.OpenFile(o, os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		for _, doc := range docsData {
			str := processing.StringConcat([]interface{}{
				doc["index"],
				"\n",
				doc["doc"],
				"\n",
			})
			f.WriteString(str)
		}
	}
}
