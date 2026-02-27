package headers

import (
	"bytes"
	"fmt"
	"slices"
	"strings"
)

const crlf = "\r\n"
var specialCharacters = []rune{'!', '#', '$', '%', '&', '\'', '*', '+', '-', '.', '^', '_', '`', '|', '~'}


type Headers map[string]string

func NewHeaders() Headers {
	return map[string]string{}
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	idx := bytes.Index(data, []byte(crlf))
	if idx == -1 {
		return 0, false, nil
	}
	if idx == 0 {
		// the empty line
		// headers are done, consume the CRLF
		return 2, true, nil
	}

	parts := bytes.SplitN(data[:idx], []byte(":"), 2)
	key := string(parts[0])

	if key != strings.TrimRight(key, " ") {
		return 0, false, fmt.Errorf("invalid header name: %s", key)
	}

	value := bytes.TrimSpace(parts[1])
	key = strings.ToLower(strings.TrimSpace(key))

	for _, r := range key {
		if (r < 'A' || r > 'Z') && (r < 'a' || r > 'z') && (r < '0' || r > '9') && !slices.Contains(specialCharacters, r){
			return 0, false, fmt.Errorf("invalid header token found: %s", key)
		}
	}

	h.Set(key, string(value))
	return idx + 2, false, nil
}

func (h Headers) Set(key, value string) {
	if elem, ok := h[key]; ok {
		h[key] = fmt.Sprintf("%s, %s", elem, value)
		return
	}

	h[key] = value
}

func (h Headers) Get(key string) (string, bool) {
	key = strings.ToLower(key)
	v, ok := h[key]
	return v, ok
}
