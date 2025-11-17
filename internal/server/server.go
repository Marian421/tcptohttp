package server

import (
	"fmt"
	"io"
	"net"

	"github.com/Marian421/tcptohttp/internal/response"
)

type Server struct {
	listener net.Listener
	closed   bool
}

func (s *Server) Close() error {
	err := s.listener.Close()
	if err != nil {
		return fmt.Errorf("error while trying to close the listener: %w", err)
	}

	s.closed = true

	return nil
}

func runConnection(conn io.ReadWriteCloser) {
	if err := response.WriteStatusLine(conn, response.StatusOk); err != nil {
		// what to do
	}
	h := response.GetDefaultHeaders(0)

	if err := response.WriteHeaders(conn, h); err != nil {
		// what to do
	}

	conn.Close()
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

		go runConnection(conn)
	}

	return nil
}

func Serve(port int) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))

	if err != nil {
		return nil, fmt.Errorf("error while trying to start the listener: %w", err)
	}

	server := &Server{
		listener: listener,
	}

	go runServer(server, listener)
	// to start listening for requests in a go routine

	return server, nil
}
