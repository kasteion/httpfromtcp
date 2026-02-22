package headers

import (
	"bytes"
	"fmt"
	"strings"
)

const clrf = "\r\n"
const colon = ":"

type Headers map[string]string

func NewHeaders() Headers {
	return Headers{}
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	if len(data) == 0 {
		return 0, false, fmt.Errorf("Bad header format")
	}

	crlfIdx := bytes.Index(data, []byte(clrf))
	if crlfIdx == -1 {
		return 0, false, nil
	}

	if crlfIdx == 0 {
		return 0, true, nil
	}

	cIdx := bytes.Index(data, []byte(colon))


	for cIdx != -1{
		key := string(data[:cIdx])
		if strings.Contains(key, " ") {
			return 0, false, fmt.Errorf("Bad header format")
		}

		value := strings.Trim(string(data[cIdx+1:crlfIdx]), " ")
		h[key] = value

		n += crlfIdx + 2

		data = data[crlfIdx+2:]
		cIdx = bytes.Index(data, []byte(colon))
		crlfIdx = bytes.Index(data, []byte(clrf))
		// fmt.Println("cIdx:", cIdx, "crlfIdx:", crlfIdx, "key:", key, "value:", value, "data:", string(data))
	}
	
	return n, false, nil
}