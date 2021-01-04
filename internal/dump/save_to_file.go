/*
Package dump provides functions for saving indexes.
*/
package dump

import (
	"bytes"
	"log"
	"os"
)

// saveDocToFile is a function that saves the passed document data a file.
func saveDocToFile(o string, docsData []map[string]string, isFirst bool) {
	if isFirst {
		f, err := os.OpenFile(o, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0666)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		var buf bytes.Buffer
		for _, doc := range docsData {
			buf.WriteString(doc["index"])
			buf.WriteString("\n")
			buf.WriteString(doc["doc"])
			buf.WriteString("\n")
		}
		f.WriteString(buf.String())
	} else {
		f, err := os.OpenFile(o, os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		var buf bytes.Buffer
		for _, doc := range docsData {
			buf.WriteString(doc["index"])
			buf.WriteString("\n")
			buf.WriteString(doc["doc"])
			buf.WriteString("\n")
		}
		f.WriteString(buf.String())
	}
}