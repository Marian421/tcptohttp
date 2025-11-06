package request

import (
	"bytes"
	"fmt"
	"io"
	"strings"
)

type Request struct {
	RequestLine RequestLine
	state       parserState
}

type parserState string

const (
	StateInit parserState = "init"
	StateDone parserState = "done"
)

func newRequest() *Request {
	return &Request{
		state: StateInit,
	}
}

// Checks if a request is done, return true if is done and false otherwise
func (r *Request) done() bool {
	return r.state == StateDone
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

var SEPARATOR = []byte("\r\n")
var ERROR_MALFORMED_START_LINE = fmt.Errorf("malformed start line")
var ERROR_UNAVAILABLE_VERSION = fmt.Errorf("unavailable http version, required 1.1")
var ERROR_MALFORMED_HTTP_FORMAT = fmt.Errorf("should be HTTP/1.1, the separator should be '/'")

// Checks that the http version is equal to 1.1
func (r *RequestLine) ValidHttp() bool {
	if r.HttpVersion == "1.1" {
		return true
	} else {
		return false
	}
}

// takes a slice of bytes and tries to parse the request line
// if succesful, returns a pointer to the request line and the number of bytes it read
func parseRequestLine(b []byte) (*RequestLine, int, error) {
	//NOTE: If it can't find a \r\n it should return 0, err
	idx := bytes.Index(b, SEPARATOR)
	if idx == -1 {
		return nil, 0, nil
	}

	line := b[:idx]

	// split string
	parts := bytes.Split(line, []byte(" "))

	// check number of items
	if len(parts) != 3 {
		return nil, idx + len(SEPARATOR), ERROR_MALFORMED_START_LINE
	}

	// split Http/1.1 in two parts
	versionParts := strings.SplitN(string(parts[2]), "/", 2)
	if len(versionParts) != 2 {
		return nil, idx + len(SEPARATOR), ERROR_MALFORMED_HTTP_FORMAT
	}

	// assign parts
	rl := &RequestLine{
		Method:        string(parts[0]),
		RequestTarget: string(parts[1]),
		HttpVersion:   versionParts[1],
	}

	// validate RequestLine items
	if !rl.ValidHttp() {
		return nil, idx + len(SEPARATOR), ERROR_UNAVAILABLE_VERSION
	}

	return rl, idx + len(SEPARATOR), nil
}

// takes a number of bytes and tries to read them
// returns the number of bytes it read
func (r *Request) parse(data []byte) (int, error) {

	switch r.state {
	case StateInit:
		rl, n, err := parseRequestLine(data)
		if n > 0 {
			r.state = StateDone
			if rl == nil {
				return n, err
			}
			r.RequestLine = *rl
			return n, err
		}
		if n == 0 {
			return n, nil
		}

		if err != nil {
			return n, err
		}

	case StateDone:
		return 0, nil
	}

	return 0, nil
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	request := newRequest()

	// NOTE: Works for now because we don't take the content currently
	buffer := make([]byte, 1024)

	// reading small chunks of data
	readPerCycle := 8

	// cummulates bytes till it parses the first line from the reader
	var working []byte

	for !request.done() {
		n, err := reader.Read(buffer[:readPerCycle])

		if n > 0 {
			working = append(working, buffer[:n]...)
		}

		// TODO: What to do here?
		if err != nil {
			return nil, err
		}

		// if it can't read a line returns 0 bytes read
		readN, err := request.parse(working)
		if readN == 0 {
			continue
		}
		if err != nil {
			return nil, err
		}

		working = working[readN:]

	}

	return request, nil
}
