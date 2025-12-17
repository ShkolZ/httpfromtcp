package server

import (
	"fmt"
	"io"

	"github.com/ShkolZ/httpfromtcp/internal/request"
	"github.com/ShkolZ/httpfromtcp/internal/response"
)

type Mux struct {
	handlers map[string]Handler
}

func NewMux() (*Mux, error) {
	mux := &Mux{
		handlers: make(map[string]Handler),
	}

	return mux, nil
}

func notFound(w io.Writer, req *request.Request) {
	headers := response.GetDefaultHeaders(0)
	response.WriteHeaders(w, headers)
	response.WriteStatusLine(w, 404)
}

func (m *Mux) Register(path string, handler Handler) {
	m.handlers[path] = handler
}

func (m *Mux) Handle(w io.Writer, req *request.Request) {
	path := req.RequestLine.RequestTarget
	fmt.Println(path)
	h, ok := m.handlers[path]

	if !ok {
		notFound(w, req)
	}

	h(w, req)

}
