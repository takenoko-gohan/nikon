/*
Package dump provides functions for saving indexes.
*/
package dump

import (
	"bytes"
	"log"
	"os"

	"github.com/takenoko-gohan/nikon/internal/processing"
)

const prefixFirstMsg = "saved "
const prefixLastMsg = " documents to a file."

// saveDocToFile is a function that saves the passed document data a file.
func saveDocToFile(o string, in <-chan []map[string]string) error {
	var f *os.File

	_, err := os.Stat(o)
	if os.IsNotExist(err) {
		f, err = os.OpenFile(o, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0666)
		if err != nil {
			return err
		}
		defer f.Close()
	} else {
		f, err = os.OpenFile(o, os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return err
		}
		defer f.Close()
	}

	for docs := range in {
		var cnt int
		var buf bytes.Buffer

		for _, doc := range docs {
			cnt++
			buf.WriteString(doc["index"])
			buf.WriteString("\n")
			buf.WriteString(doc["doc"])
			buf.WriteString("\n")
		}
		_, err := f.WriteString(buf.String())
		if err != nil {
			return err
		}

		msg := processing.StringConcat([]interface{}{
			prefixFirstMsg,
			cnt,
			prefixLastMsg,
		})
		log.Println(msg)
	}

	return nil
}
