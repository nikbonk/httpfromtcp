package server

import (
	"log"
	"net"
	"strconv"
	"sync/atomic"
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
	response := []byte("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\n\r\nHello World!\n")
	_, err := conn.Write(response)
	if err != nil {
		log.Println(err)
		_ = conn.Close()
		return
	}
	_ = conn.Close()
}
