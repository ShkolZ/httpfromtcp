package response

import (
	"fmt"
	"io"

	"github.com/ShkolZ/httpfromtcp/internal/headers"
)

type StatusCode int

const (
	OK                  StatusCode = 200
	BadRequest          StatusCode = 400
	InternalServerError StatusCode = 500
	NotFound            StatusCode = 404
)

func WriteStatusLine(w io.Writer, sc StatusCode) error {
	var sLine string
	switch sc {
	case OK:
		sLine = "HTTP/1.1 200 OK"
	case BadRequest:
		sLine = "HTTP/1.1 400 Bad Request"
	case InternalServerError:
		sLine = "HTTP/1.1 500 Internal Server Error"
	case NotFound:
		sLine = "HTTP/1.1 404 Not Found"
	default:
		sLine = "HTTP/1.1 403 Idk"
	}
	_, err := w.Write([]byte(fmt.Sprintf("%s\r\n", sLine)))
	if err != nil {
		return err
	}
	fmt.Println(sLine)
	return nil
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	h := headers.NewHeaders()
	h.Set("content-type", "text/plain")
	h.Set("connection", "close")
	h.Set("content-length", fmt.Sprintf("%v", contentLen))

	return h
}

func WriteHeaders(w io.Writer, h headers.Headers) error {
	for key, val := range h {
		_, err := w.Write([]byte(fmt.Sprintf("%s: %s\r\n", key, val)))
		if err != nil {
			return err
		}
	}
	w.Write([]byte("\r\n"))
	return nil
}

func WriteBody(w io.Writer, data []byte) {
	headers := GetDefaultHeaders(len(data))
	WriteStatusLine(w, 200)
	WriteHeaders(w, headers)
	w.Write(data)
}
