package request

import (
	"errors"
	"io"
	"regexp"
	"strings"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

var IsCapitalizedLetter = regexp.MustCompile(`^[A-Z]+$`).MatchString

func RequestFromReader(reader io.Reader) (*Request, error) {
	rb, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	rs := string(rb)

	lines := strings.Split(rs, "\r\n")
	if len(lines) == 0 {
		return nil, errors.New("empty reader")
	}

	requestLine, err := parseRequestLine(lines[0])
	if err != nil {
		return nil, err
	}

	parsedRequest := &Request{
		RequestLine: *requestLine,
	}

	return parsedRequest, nil
}

func parseRequestLine(rs string) (*RequestLine, error) {
	parts := strings.Split(rs, " ")
	if len(parts) != 3 {
		return nil, errors.New("invalid number of parts in request line")
	}

	method := parts[0]
	if !IsCapitalizedLetter(method) {
		return nil, errors.New("invalid HTTP method")
	}

	protocolParts := strings.Split(parts[2], "/")
	if len(protocolParts) != 2 {
		return nil, errors.New("invalid protocol")
	}

	if protocolParts[0] != "HTTP" {
		return nil, errors.New("invalid protocol")
	}

	httpVersion := protocolParts[1]
	if httpVersion != "1.1" {
		return nil, errors.New("invalid protocol version")
	}

	requestTarget := parts[1]

	return &RequestLine{
		HttpVersion:   httpVersion,
		RequestTarget: requestTarget,
		Method:        method,
	}, nil
}
