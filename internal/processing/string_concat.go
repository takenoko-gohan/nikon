package processing

import (
	"bytes"
	"strconv"
)

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
