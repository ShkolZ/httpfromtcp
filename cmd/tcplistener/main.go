package main

import (
	"fmt"
	"log"
	"net"

	"github.com/ShkolZ/httpfromtcp/internal/request"
)

func main() {
	ln, err := net.Listen("tcp", ":42069")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Server is listening on port 42069...")
	defer ln.Close()
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println("Connection was accepted")

		req, err := request.RequestFromReader(conn)
		if err != nil {
			log.Fatalln(err)
		}

		fmt.Printf(`
Request line:
- Method: %v
- Target: %v
- Version: %v
`, req.RequestLine.Method,
			req.RequestLine.RequestTarget,
			req.RequestLine.HttpVersion)
		fmt.Println("Headers:")
		for k, v := range req.Headers {
			fmt.Printf("- %v: %v\n", k, v)
		}
		fmt.Println("Body:")
		fmt.Println(string(req.Body))

		conn.Close()
		fmt.Println("Connection closed")
	}
}
