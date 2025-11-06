package request

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"
)

type Request struct {
	RequestLine RequestLine
	state       parserState
}

func newRequest() *Request {
	return &Request{
		state: StateInit,
	}
}

type parserState string

const (
	StateInit parserState = "init"
	StateDone parserState = "done"
)

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

var SEPARATOR = "\r\n"
var ERROR_MALFORMED_START_LINE = fmt.Errorf("malformed start line")
var ERROR_UNAVAILABLE_VERSION = fmt.Errorf("unavailable http version")

func (r *RequestLine) ValidHttp() bool {
	if r.HttpVersion == "1.1" {
		return true
	} else {
		return false
	}
}

// splits the HTTP message and returns the request line and the rest of the message
func parseRequestLine(b string) (*RequestLine, string, error) {
	//get line
	//TODO: return how many the number of bytes it consumed
	//NOTE: If it can't find a \r\n it should return 0, err
	idx := strings.Index(b, SEPARATOR)
	if idx == -1 {
		return nil, b, nil
	}

	startLine := b[:idx]
	restOfMsg := b[idx+len(SEPARATOR):]

	// split string
	parts := strings.Split(startLine, " ")

	// check number of items
	if len(parts) != 3 {
		return nil, restOfMsg, ERROR_MALFORMED_START_LINE
	}

	// assign parts
	rl := &RequestLine{
		Method:        parts[0],
		RequestTarget: parts[1],
		HttpVersion:   strings.Split(parts[2], "/")[1],
	}

	// validate RequestLine items
	if !rl.ValidHttp() {
		return nil, restOfMsg, ERROR_UNAVAILABLE_VERSION
	}

	return rl, restOfMsg, nil

}

func (r *Request) parse(data []byte) (int, error) {
	read := 0

outer:
	switch r.state {
	case StateInit:
	}
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	request := newRequest()
	// make a buffer
	buffer := make([]byte, 8)
	bufIdx := 0
	for {
		n, err := reader.Read(buffer[bufIdx:])
		// TODO: What to do here?
		if err != nil {
			return nil, err
		}

		readN, err := r.parse()
	}

	return &Request{
		RequestLine: *rl,
	}, err
}
