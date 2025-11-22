package main

import (
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Marian421/tcptohttp/internal/request"
	"github.com/Marian421/tcptohttp/internal/response"
	"github.com/Marian421/tcptohttp/internal/server"
)

const port = 42069

func main() {
	server, err := server.Serve(port, func(w io.Writer, req *request.Request) *server.HandlerError {
		handlerErr := new(server.HandlerError)

		if req.RequestLine.RequestTarget == "/myproblem" {
			handlerErr.Status = response.StatusInternalError
			handlerErr.Message = "Woopsie, my bad!\n"
			return handlerErr
		} else if req.RequestLine.RequestTarget == "/yourproblem" {
			handlerErr.Status = response.StatusInternalError
			handlerErr.Message = "Your bad\n"
			return handlerErr
		} else {
			w.Write([]byte("All good frfr\n"))
		}

		return nil
	})
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}
