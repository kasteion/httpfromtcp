package response

import (
	"fmt"
	"io"

	"github.com/kasteion/httpfromtcp/internal/headers"
)

type StatusCode int

const (
	StatusCodeOK StatusCode = 200
	StatusCodeBadRequest StatusCode = 400
	StatusCodeInternalServerError StatusCode = 500
)

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	statusLine := ""
	switch statusCode {
	case StatusCodeOK:
		statusLine = "HTTP/1.1 200 OK\r\n"
	case StatusCodeBadRequest:
		statusLine = "HTTP/1.1 400 Bad Request\r\n"
	case StatusCodeInternalServerError:
		statusLine = "HTTP/1.1 500 Internal Server Error\r\n"
	default:
		statusLine = fmt.Sprintf("HTTP/1.1 %d \r\n", statusCode)
	}
	
	_, err := w.Write([]byte(statusLine))
	if err != nil {
		return err
	}
	return nil
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	h := headers.NewHeaders()
	h.Set("Content-Length", fmt.Sprint(contentLen))
	h.Set("Connection", "close")
	h.Set("Content-Type", "text/plain")
	return h
}

func WriteHeaders(w io.Writer, headers headers.Headers) error {
	for k := range headers {
		v := headers[k]
		fieldLine := fmt.Sprintf("%s: %v\r\n", k, v)
		_, err := w.Write([]byte(fieldLine))
		if err != nil {
			return err
		}
	}
	w.Write([]byte("\r\n"))
	return nil
}
