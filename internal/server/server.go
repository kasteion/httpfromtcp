package server

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"sync/atomic"

	"github.com/kasteion/httpfromtcp/internal/request"
	"github.com/kasteion/httpfromtcp/internal/response"
)

type Server struct {
	closed atomic.Bool
	listener net.Listener
	handler Handler
}

type HandlerError struct {
	StatusCode response.StatusCode
	Message string
}

type Handler func(w io.Writer, req *request.Request) *HandlerError

func Serve(port int, handler Handler) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}
	server := &Server{ 
		listener: listener,
		handler: handler,
	}
	go server.listen()
	return server, nil
}

func (s *Server) Close() error {
	s.closed.Store(true)
	if s.listener != nil {
		return s.listener.Close()
	}
	return nil
}

func (s *Server) listen() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			if s.closed.Load() {
				return
			}
			log.Printf("could not accept connection: %s", err)
			continue
		}
		go s.handle(conn)
	}
}

func (err HandlerError) Write(w io.Writer) {
	response.WriteStatusLine(w, err.StatusCode)
	h := response.GetDefaultHeaders(len(err.Message))
	response.WriteHeaders(w, h)
	w.Write([]byte(err.Message))
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()
	req, err := request.RequestFromReader(conn)
	if err != nil {
		log.Printf("could not read reqeust: %s", err)
		return
	}

	var buf bytes.Buffer
	hError := s.handler(&buf, req)
	if hError != nil {
		hError.Write(conn)
		return
	}

	response.WriteStatusLine(conn, response.StatusCodeOK)
	headers := response.GetDefaultHeaders(buf.Len())
	response.WriteHeaders(conn, headers)
	conn.Write(buf.Bytes())
}
