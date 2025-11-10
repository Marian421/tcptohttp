package main

import (
	"fmt"
	"github.com/Marian421/tcptohttp/internal/request"
	"log"
	"net"
)

func main() {
	listener, err := net.Listen("tcp", ":42069")
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()
	fmt.Println("Listening")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}

		r, err := request.RequestFromReader(conn)
		if err != nil {
			log.Fatal("error", "error", err)
		}

		fmt.Printf("Request line:\n")
		fmt.Printf("- Method: %s\n", r.RequestLine.Method)
		fmt.Printf("- Target: %s\n", r.RequestLine.RequestTarget)
		fmt.Printf("- Version: %s\n", r.RequestLine.HttpVersion)
		fmt.Println("Headers:")
		for key, value := range r.Headers {
			fmt.Printf("%s: %s\n", key, value)
		}
		fmt.Println(string(r.Body))
	}
}
