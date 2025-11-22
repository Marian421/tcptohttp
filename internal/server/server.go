package server

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net"

	"github.com/Marian421/tcptohttp/internal/request"
	"github.com/Marian421/tcptohttp/internal/response"
)

type Server struct {
	listener net.Listener
	closed   bool
	handler  Handler
}

type HandlerError struct {
	Status  response.StatusCode
	Message string
}

type Handler func(w io.Writer, req *request.Request) *HandlerError

func (s *Server) Close() error {
	err := s.listener.Close()
	if err != nil {
		return fmt.Errorf("error while trying to close the listener: %w", err)
	}

	s.closed = true

	return nil
}

func runConnection(conn io.ReadWriteCloser, s *Server) {
	defer conn.Close()

	r, err := request.RequestFromReader(conn)
	if err != nil {
		return
	}

	var buf bytes.Buffer
	handlerErr := s.handler(&buf, r)

	// Wrap the connection to flush all writes immediately
	bw := bufio.NewWriter(conn)

	if handlerErr != nil {
		fmt.Printf("message: %s", handlerErr.Message)
		fmt.Printf("status: %d", handlerErr.Status)

		response.WriteStatusLine(bw, handlerErr.Status)
		response.WriteHeaders(bw, response.GetDefaultHeaders(len(handlerErr.Message)))
		bw.Write([]byte(handlerErr.Message))

		bw.Flush() // <- THIS ensures all bytes go to the socket immediately
		return
	}

	response.WriteStatusLine(bw, response.StatusOk)
	response.WriteHeaders(bw, response.GetDefaultHeaders(buf.Len()))
	bw.Write(buf.Bytes())
	bw.Flush() // flush normal response
}

func WriteErrors(w io.Writer, err *HandlerError) {
	fmt.Fprintf(w, "Error %d: %s\n", err.Status, err.Message)
}

func runServer(s *Server, listener net.Listener) error {

	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}

		if s.closed == true {
			return nil
		}

		go runConnection(conn, s)
	}

	return nil
}

func Serve(port int, handler Handler) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))

	if err != nil {
		return nil, fmt.Errorf("error while trying to start the listener: %w", err)
	}

	server := &Server{
		listener: listener,
		handler:  handler,
	}

	go runServer(server, listener)
	// to start listening for requests in a go routine

	return server, nil
}
