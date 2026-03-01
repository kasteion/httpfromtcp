package main

import (
	"fmt"
	"log"
	"net"

	"github.com/kasteion/httpfromtcp/internal/request"
)

const port = ":42069"

func main() {
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("could not listen on %s - %s\n", port, err)
	}
	defer listener.Close()

	fmt.Printf("Listening on %s\n", port)
	fmt.Println("===================")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalf("could not accept connection: %s", err)
		}
		fmt.Println("Connection accepted from", conn.RemoteAddr())

		req, err := request.RequestFromReader(conn)
		if err != nil {
			log.Fatalf("error parsing request: %s\n", err.Error())
		}
		fmt.Println("Request line:")
		fmt.Println("- Method:", req.RequestLine.Method)
		fmt.Println("- Target:", req.RequestLine.RequestTarget)
		fmt.Println("- Version:", req.RequestLine.HttpVersion)

		fmt.Println("Headers:")
		for key := range(req.Headers) {
			fmt.Printf("- %s: %s\n", key, req.Headers[key])
		}

		fmt.Println("Body:")
		fmt.Println(string(req.Body))

		fmt.Println("Connection to", conn.RemoteAddr(), "closed")
	}
}
