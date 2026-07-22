package server

import (
	"log"
	"net"
	"strconv"
	"sync/atomic"

	"github.com/nikbonk/httpfromtcp/internal/response"
)

type Server struct {
	listener net.Listener
	closed   atomic.Bool
}

func Serve(port int) (*Server, error) {
	listenerPort := ":" + strconv.Itoa(port)
	listener, err := net.Listen("tcp", listenerPort)
	if err != nil {
		return nil, err
	}
	server := &Server{
		listener: listener,
		closed:   atomic.Bool{},
	}
	go server.listen()

	return server, nil
}

func (s *Server) Close() error {
	if s.closed.CompareAndSwap(false, true) {
		return s.listener.Close()
	}
	return nil
}

func (s *Server) listen() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			return
		}
		go s.handle(conn)
	}
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()
	err := response.WriteStatusLine(conn, response.StatusCodeOK)
	if err != nil {
		log.Println(err)
		return
	}
	contentLen := 0
	headers := response.GetDefaultHeaders(contentLen)
	err = response.WriteHeaders(conn, headers)
	if err != nil {
		log.Println(err)
		return
	}

}
