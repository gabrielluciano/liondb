package server

import (
	"fmt"
	"io"
	"net"
)

type MessageHandler func(message string) []byte

type Server struct {
	port           string
	messageHandler MessageHandler
}

func New(port string) *Server {
	return &Server{port: port}
}

func (s *Server) SetMessageHandler(fn MessageHandler) {
	s.messageHandler = fn
}

func (s *Server) Listen() {
	if s.messageHandler == nil {
		panic("Message handler not defined!")
	}
	ln, err := net.Listen("tcp", ":"+s.port)
	if err != nil {
		panic(err)
	}
	fmt.Println("Server started on port 7123")

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Printf("Error accepting connection: %v", err)
		}
		go s.handleConnection(conn)
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	buffer := make([]byte, 512)
	for {
		n, err := conn.Read(buffer)
		if err == io.EOF {
			break
		}
		response := s.messageHandler(string(buffer[:n]))
		conn.Write(append(response, '\n'))
	}
}
