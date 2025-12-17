package request

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/ShkolZ/httpfromtcp/internal/headers"
)

const (
	GET        = "GET"
	POST       = "POST"
	PUT        = "PUT"
	DELETE     = "DELETE"
	bufferSize = 1024
)

type ParserState string

const (
	StateInit ParserState = "init"
	StateDone ParserState = "done"
)

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

type Request struct {
	RequestLine *RequestLine
	Headers     headers.Headers
	Body        []byte
	parserState ParserState
	headerDone  bool
}

var reqSeparator = "\r\n"
var Error_Badly_Formatted = fmt.Errorf("not valid format of request-line")

func (r *Request) isDone() bool {
	return r.parserState == StateDone
}

func (r *Request) finish() {
	r.parserState = StateDone
}

func (r *Request) parse(data []byte) (int, error) {
	read := 0
	bodyDone := false

outer:
	for {
		switch r.parserState {
		case StateInit:
			if r.RequestLine == nil {
				reqLine, n, err := parseRequestLine(data)
				if err != nil {
					return read, err
				}

				if n == 0 {
					break outer
				}

				r.RequestLine = reqLine
				read += n
			}

			if !r.headerDone {
				n, done, err := r.Headers.Parse(data[read:])
				if err != nil {
					return 0, err
				}

				if n == 0 {
					break outer
				}

				read += n
				r.headerDone = done
			}

			val, ok := r.Headers.Get("content-length")
			if ok && r.headerDone {

				intVal, _ := strconv.Atoi(val)
				body, n, _ := parseBody(data[read:], intVal)

				if n == 0 {
					break outer
				}

				r.Body = body
				bodyDone = true
				read += n
			} else {
				bodyDone = true
			}
			if r.headerDone && bodyDone {
				r.finish()
			}
		case StateDone:
			break outer
		}
	}
	return read, nil
}

func RequestFromReader(r io.Reader) (*Request, error) {
	buf := make([]byte, bufferSize)
	bufLen := 0
	req := &Request{
		parserState: StateInit,
		Headers:     headers.Headers{},
	}

	for !req.isDone() {
		n, err := r.Read(buf[bufLen:])
		if err != nil {
			return req, nil
		}

		if n == 0 {
			return req, fmt.Errorf("some error")
		}

		bufLen += n

		readN, err := req.parse(buf[:bufLen])
		if err != nil {
			return nil, err
		}

		copy(buf, buf[readN:bufLen])
		bufLen -= readN
		if len(buf) == bufLen {
			buf = append(buf, 0)
		}
	}

	return req, nil
}

func parseRequestLine(b []byte) (*RequestLine, int, error) {

	sepIdx := bytes.Index(b, []byte(reqSeparator))
	if sepIdx == -1 {
		return nil, 0, nil
	}
	rlStr := b[:sepIdx]

	rlSlice := strings.Split(string(rlStr), " ")
	if len(rlSlice) != 3 {
		return &RequestLine{}, 0, Error_Badly_Formatted
	}
	met := rlSlice[0]
	rp := rlSlice[1]
	ver := strings.Split(rlSlice[2], "/")

	for _, c := range met {
		if !(c > 64 && c < 91) {
			return &RequestLine{}, 0, fmt.Errorf("method should be uppercase")
		}
	}

	if !(met == GET || met == POST || met == PUT || met == DELETE) {
		return &RequestLine{}, 0, fmt.Errorf("not valid method")
	}
	if rp[0] != '/' {
		return &RequestLine{}, 0, fmt.Errorf("not valid resource path")
	}

	if !(ver[0] == "HTTP" && ver[1] == "1.1") {
		return &RequestLine{}, 0, fmt.Errorf("doesn't support this version: %v", ver[1])
	}

	reqLine := &RequestLine{
		HttpVersion:   ver[1],
		RequestTarget: rp,
		Method:        met,
	}

	return reqLine, sepIdx + 2, nil
}

func parseBody(data []byte, bl int) ([]byte, int, bool) {
	if len(data) == bl {
		body := data
		return body, bl, true
	}
	return []byte{}, 0, false
}
