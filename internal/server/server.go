package server

import (
	"fmt"
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

type Handler func(w *response.Writer, req *request.Request)

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

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()
	req, err := request.RequestFromReader(conn)
	if err != nil {
		log.Printf("could not read reqeust: %s", err)
		return
	}

	w := response.NewWriter(conn)
	s.handler(w, req)
}
