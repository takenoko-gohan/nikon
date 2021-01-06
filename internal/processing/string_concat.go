/*
Package processing provides functions to 
process the values passed in.
*/
package processing

import (
	"bytes"
	"strconv"
)

/*
StringConcat is a function that concatenates the
passed values and returns them as a string.
*/
func StringConcat(values []interface{}) string {
	var buf bytes.Buffer

	for _, value := range values {
		switch value.(type) {
		case string:
			buf.WriteString(value.(string))
		case int:
			s := strconv.Itoa(value.(int))
			buf.WriteString(s)
		}
	}

	return buf.String()
}
