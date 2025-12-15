package headers

import (
	"bytes"
	"fmt"
	"strings"
)

const separator = "\r\n"

var white_space_error = fmt.Errorf("unallowed white space")
var special_character_error = fmt.Errorf("unallowed character")

type Headers map[string]string

func (h Headers) Set(key, val string) {
	key = strings.ToLower(key)
	val = strings.TrimSpace(val)
	prev, ok := h[key]
	if ok {
		h[key] = fmt.Sprintf("%v, %v", prev, val)
		return
	}
	h[key] = val
}

func (h Headers) Get(key string) (string, bool) {
	val, ok := h[strings.ToLower(key)]
	return val, ok
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	headerLen := bytes.Index(data, []byte(separator))
	if headerLen == -1 {
		return 0, false, nil
	}

	headerLine := string(data[:headerLen])

	if data[0] == 13 && data[1] == 10 {
		return headerLen + 2, true, nil
	}

	trHl := strings.TrimSpace(headerLine)
	parts := strings.Split(trHl, ": ")
	if strings.Contains(parts[0], " ") {
		return 0, false, white_space_error
	}

	if strings.Contains(parts[1][1:], " ") {
		return 0, false, white_space_error
	}

	for _, c := range parts[0] {
		if (c > 32 && c < 45) || (c > 45 && c < 65) {
			return 0, false, special_character_error
		}
	}

	h.Set(parts[0], parts[1])
	return headerLen + 2, false, nil

}

func NewHeaders() Headers {
	return Headers{}
}
