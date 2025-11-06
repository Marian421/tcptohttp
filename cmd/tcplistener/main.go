package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
)

func getLinesChannel(c io.ReadCloser) <-chan string {
	ch := make(chan string, 1)

	go func() {
		defer close(ch)
		defer c.Close()

		buffer := make([]byte, 8)
		leftover := []byte{}

		for {
			n, err := c.Read(buffer) // number of bytes read
			if err != nil {          // checks for error
				if err == io.EOF {
					break
				}
				log.Fatal(err)
			}

			chunk := append(leftover, buffer[:n]...) // adds the bytes that have been read to a chunk

			if i := bytes.IndexByte(chunk, '\n'); i == -1 { // checks for the end of line
				leftover = chunk
			} else {
				ch <- string(chunk[:i])
				leftover = chunk[i+1:]
			}
		}
	}()

	return ch
}

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

		for line := range getLinesChannel(conn) {
			fmt.Printf("%s\n", line)
		}

	}

}
