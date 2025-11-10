package request

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"strings"
	"testing"
)

type chunkReader struct {
	data            string
	numBytesPerRead int
	pos             int
}

func (cr *chunkReader) Read(p []byte) (n int, err error) {
	if cr.pos >= len(cr.data) {
		return 0, io.EOF
	}
	endIndex := cr.pos + cr.numBytesPerRead
	if endIndex > len(cr.data) {
		endIndex = len(cr.data)
	}
	n = copy(p, cr.data[cr.pos:endIndex])
	cr.pos += n
	return n, nil
}

func assertRequestLine(t *testing.T, r *Request, method, target, version string) {
	t.Helper()
	assert.Equal(t, method, r.RequestLine.Method)
	assert.Equal(t, target, r.RequestLine.RequestTarget)
	assert.Equal(t, version, r.RequestLine.HttpVersion)
}

func TestRequestFromReader(t *testing.T) {

	t.Run("Parses full request with body", func(t *testing.T) {
		r, err := RequestFromReader(strings.NewReader(
			"GET / HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n{background-color: blue}",
		))
		require.NoError(t, err)
		require.NotNil(t, r)

		assertRequestLine(t, r, "GET", "/", "1.1")
		assert.Equal(t, "{background-color: blue}", string(r.Body))
	})

	t.Run("Content-Length has the correct value", func(t *testing.T) {
		r, err := RequestFromReader(strings.NewReader(
			"GET / HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\nContent-Length: 24\r\n\r\n{background-color: blue}",
		))
		require.NoError(t, err)
		require.NotNil(t, r)

		assertRequestLine(t, r, "GET", "/", "1.1")
		assert.Equal(t, "{background-color: blue}", string(r.Body))
	})

	t.Run("Content-Length has the correct value", func(t *testing.T) {
		r, err := RequestFromReader(strings.NewReader(
			"GET / HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\nContent-Length: 24\r\n\r\n{background-color: blue}",
		))
		require.NoError(t, err)
		require.NotNil(t, r)

		assertRequestLine(t, r, "GET", "/", "1.1")
		assert.Equal(t, "{background-color: blue}", string(r.Body))
	})

	t.Run("Content-Length has a wrong value", func(t *testing.T) {
		_, err := RequestFromReader(strings.NewReader(
			"GET / HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\nContent-Length: 23\r\n\r\n{background-color: blue}",
		))
		require.Error(t, err)
	})

	t.Run("Parses request without body", func(t *testing.T) {
		r, err := RequestFromReader(strings.NewReader(
			"GET / HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
		))
		require.NoError(t, err)
		require.NotNil(t, r)

		assertRequestLine(t, r, "GET", "/", "1.1")
		assert.Equal(t, "curl/7.81.0", r.Headers.Get("User-Agent"))
		assert.Equal(t, "*/*", r.Headers.Get("Accept"))
	})

	t.Run("Fails if no CRLF before body", func(t *testing.T) {
		_, err := RequestFromReader(strings.NewReader(
			"GET / HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n",
		))
		require.Error(t, err)
	})

	t.Run("Parses request with path", func(t *testing.T) {
		r, err := RequestFromReader(strings.NewReader(
			"GET /coffee HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
		))
		require.NoError(t, err)
		require.NotNil(t, r)

		assertRequestLine(t, r, "GET", "/coffee", "1.1")
	})

	t.Run("Fails when request line is malformed", func(t *testing.T) {
		_, err := RequestFromReader(strings.NewReader(
			"/coffee HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
		))
		require.Error(t, err)
	})

	t.Run("Reads correctly with chunked Reader", func(t *testing.T) {
		t.Run("Simple GET", func(t *testing.T) {
			reader := &chunkReader{
				data:            "GET / HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
				numBytesPerRead: 3,
			}
			r, err := RequestFromReader(reader)
			require.NoError(t, err)
			require.NotNil(t, r)
			assertRequestLine(t, r, "GET", "/", "1.1")
		})

		t.Run("GET with path", func(t *testing.T) {
			reader := &chunkReader{
				data:            "GET /coffee HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
				numBytesPerRead: 1,
			}
			r, err := RequestFromReader(reader)
			require.NoError(t, err)
			require.NotNil(t, r)
			assertRequestLine(t, r, "GET", "/coffee", "1.1")
		})
	})

	t.Run("Rejects unsupported HTTP versions", func(t *testing.T) {
		_, err := RequestFromReader(strings.NewReader("GET / HTTP/2.0\r\n"))
		require.Error(t, err)
	})
}

