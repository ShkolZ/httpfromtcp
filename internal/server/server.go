package server

import (
	"bytes"
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
	handler   Handler
}

type HandlerError struct {
	StatusCode response.StatusCode
	Msg        string
}

type Handler func(w io.Writer, req *request.Request) *HandlerError

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()
	req, err := request.RequestFromReader(conn)
	fmt.Println(req)
	if err != nil {
		log.Fatalln()
	}
	buf := &bytes.Buffer{}

	hErr := s.handler(buf, req)
	if hErr != nil {
		h := response.GetDefaultHeaders(len(hErr.Msg))
		response.WriteStatusLine(conn, 400)
		response.WriteHeaders(conn, h)

		return
	}

	headers := response.GetDefaultHeaders(int(buf.Len()))
	response.WriteStatusLine(conn, 200)
	response.WriteHeaders(conn, headers)
	conn.Write(buf.Bytes())
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

func Serve(port int, handler Handler) (*Server, error) {
	ln, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		return nil, err
	}

	ser := &Server{
		port:      port,
		isServing: true,
		listener:  ln,
		handler:   handler,
	}

	go ser.listen()

	return ser, nil
}
