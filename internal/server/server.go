package server

import (
	"fmt"
	"log"
	"net"
	"sync/atomic"

	"github.com/kasteion/httpfromtcp/internal/response"
)

type Server struct {
	closed atomic.Bool
	listener net.Listener
}

func Serve(port int) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}
	server := &Server{ listener: listener }
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

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()
	err := response.WriteStatusLine(conn, response.StatusCodeOK)
	if err != nil {
		log.Printf("could not write status line: %s", err)
	}
	headers := response.GetDefaultHeaders(0)
	err = response.WriteHeaders(conn, headers)
	if err != nil {
		log.Printf("could not write headers: %s", err)
	}
}
