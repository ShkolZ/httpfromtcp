package main

import (
	"fmt"
	"io"

	"github.com/ShkolZ/httpfromtcp/internal/request"
	"github.com/ShkolZ/httpfromtcp/internal/response"
)

func HelloWorldHandler(w io.Writer, req *request.Request) {
	data := []byte("Hello World!")
	fmt.Println("tyt")
	response.WriteBody(w, data)
}
