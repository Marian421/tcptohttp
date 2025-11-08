package headers

import (
	"bytes"
	"fmt"
	"strings"
)

type Headers map[string]string

func newHeaders() Headers {
	return make(Headers)
}

func (h Headers) hasKey(key string) bool {
	hasKey := false

	if _, ok := h[key]; ok {
		hasKey = true
	}
	return hasKey
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

// check is rune is valid
func isTchar(r rune) bool {
	if r >= 'a' && r <= 'z' {
		return true
	}
	if r >= 'A' && r <= 'Z' {
		return true
	}
	if r >= '0' && r <= '9' {
		return true
	}

	switch r {
	case '!', '#', '$', '%', '&', '\'', '*', '+', '-', '.', '^',
		'_', '`', '|', '~':
		return true
	}

	return false
}

// checks is a field-name has valid characters
func isValidFieldName(s string) bool {
	if len(s) == 0 {
		return false
	}

	for _, r := range s {
		if !isTchar(r) {
			return false
		}
	}

	return true
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
			read = read + len(sep)
			break
		}

		lineEnd := read + idx

		header, value, err := parseHeader(data[read:lineEnd])
		if err != nil {
			return read, done, fmt.Errorf("error while trying to parse bytes %d-%d: %w", read, lineEnd, err)
		}

		// check for special characters in field-name
		if !isValidFieldName(header) {
			return read, done, fmt.Errorf("field-name contains invalid characters")
		}

		// field-name to lowercase
		header = strings.ToLower(header)

		// check if key already exists, if yes, format field-value to be
		// value1, value2, ..., valueN

		if h.hasKey(header) {
			value = h[header] + ", " + value
		}
		// add the field-name: field-value pair in the map
		h[header] = value

		read = lineEnd + len(sep)
	}
	return read, done, nil
}
