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

type WriteState int

const (
	writerStateStatusLine WriteState = iota
	writerStateHeaders
	writerBody
)

type Writer struct {
	writer io.Writer
	writeState WriteState
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{
		writer: w,
		writeState: writerStateStatusLine,
	}
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
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
	
	_, err := w.writer.Write([]byte(statusLine))
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

func (w *Writer) WriteHeaders(headers headers.Headers) error {
	for k := range headers {
		v := headers[k]
		fieldLine := fmt.Sprintf("%s: %v\r\n", k, v)
		_, err := w.writer.Write([]byte(fieldLine))
		if err != nil {
			return err
		}
	}
	w.writer.Write([]byte("\r\n"))
	return nil
}

func (w *Writer) WriteBody(p []byte) (int, error) {
	return w.writer.Write(p)
}

func (w *Writer) WriteChunkedBody(p []byte) (int, error) {
	return w.writer.Write(fmt.Appendf(nil, "%X\r\n%s\r\n", len(p), string(p)))
}

func (w *Writer) WriteChunkedBodyDone() (int, error) {
	return w.writer.Write([]byte("0\r\n"))
}

func (w *Writer) WriteTrailers(h headers.Headers) error {
	for t := range h {
		w.writer.Write(fmt.Appendf(nil, "%s: %s\r\n", t, h[t]))
	}
	w.writer.Write([]byte("\r\n"))
	return nil
}