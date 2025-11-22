package response

import (
	"fmt"
	"io"

	"github.com/Marian421/tcptohttp/internal/headers"
)

type StatusCode int

const (
	StatusOk            StatusCode = 200
	StatusBad           StatusCode = 400
	StatusInternalError StatusCode = 500
)

func WriteHeaders(w io.Writer, headers headers.Headers) error {
	for key, value := range headers {
		_, err := w.Write([]byte(fmt.Sprintf("%s: %s\r\n", key, value)))
		if err != nil {
			return err
		}
	}

	w.Write([]byte("\r\n"))

	return nil
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	h := headers.NewHeaders()

	h.Set("Content-Length", fmt.Sprintf("%d", contentLen))
	h.Set("Connection", "close")
	h.Set("Content-Type", "text/plain")

	return h
}

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	statusLine := ""
	switch statusCode {
	case StatusOk:
		statusLine = "HTTP/1.1 200 OK\r\n"
	case StatusBad:
		statusLine = "HTTP/1.1 400 Bad Request\r\n"
	case StatusInternalError:
		statusLine = "HTTP/1.1 500 Internal Server Error\r\n"
	default:
		statusLine = "Unknown status code\r\n"
	}

	_, err := w.Write([]byte(statusLine))
	if err != nil {
		return err
	}

	return nil
}
