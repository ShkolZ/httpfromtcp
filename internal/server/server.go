package server

import (
	"fmt"
	"io"
	"log"
	"net"
	"strconv"

	"github.com/ShkolZ/httpfromtcp/internal/request"
	"github.com/ShkolZ/httpfromtcp/internal/response"
)

type Server struct {
	port      int
	isServing bool
	listener  net.Listener
	mux       *Mux
}

type HandlerError struct {
	StatusCode response.StatusCode
	Msg        string
}

type Handler func(w io.Writer, req *request.Request)

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()
	req, err := request.RequestFromReader(conn)
	fmt.Println(req)
	if err != nil {
		log.Fatalln()
	}

	s.mux.Handle(conn, req)

}

func (s *Server) listen() {
	for {
		if s.isServing {
			conn, _ := s.listener.Accept()
			go s.handle(conn)
		}
	}
}

func (s *Server) Close() error {
	s.isServing = false
	return nil
}

func Serve(port int, m *Mux) (*Server, error) {
	ln, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		return nil, err
	}

	ser := &Server{
		port:      port,
		isServing: true,
		listener:  ln,
		mux:       m,
	}

	go ser.listen()

	return ser, nil
}
