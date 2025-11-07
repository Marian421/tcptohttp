package headers

import (
	"bytes"
	"fmt"
)

type Headers map[string]string

func newHeaders() Headers {
	return make(Headers)
}

var sep = []byte("\r\n")

// takes a request line without the separator
// returns (header, value, err)
func parseHeader(line []byte) (string, string, error) {
	fieldNameParts := bytes.SplitN(line, []byte(":"), 2)
	if len(fieldNameParts) != 2 {
		return "", "", fmt.Errorf("malformed field-line")
	}

	if bytes.Index(fieldNameParts[0], []byte(" ")) != -1 {
		return "", "", fmt.Errorf("malformed field name")
	}

	trimmedValue := bytes.Trim(fieldNameParts[1], " ")

	return string(fieldNameParts[0]), string(trimmedValue), nil
}

func (h Headers) Parse(data []byte) (int, bool, error) {
	read := 0
	done := false

	for {
		idx := bytes.Index(data[read:], sep)

		if idx == -1 {
			return read, done, nil
		}

		// Empty line - the end of header
		if idx == 0 {
			done = true
			break
		}

		lineEnd := read + idx

		header, value, err := parseHeader(data[read:lineEnd])
		if err != nil {
			return read, done, fmt.Errorf("error while trying to parse bytes %d-%d: %w", read, lineEnd, err)
		}

		// add the pair in the map
		h[header] = value

		read = lineEnd + len(sep)
	}
	return read, done, nil
}
