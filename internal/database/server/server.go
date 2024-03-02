package server

import (
	"fmt"
	"io"
	"net"
)

type MessageHandler func(message string) string

var messageHandler MessageHandler

func SetMessageHandler(fn MessageHandler) {
	messageHandler = fn
}

func Listen() {
	if messageHandler == nil {
		panic("Message handler not defined!")
	}
	ln, err := net.Listen("tcp", ":7123")
	if err != nil {
		panic(err)
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Printf("Error accepting connection: %v", err)
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	buffer := make([]byte, 512)
	for {
		n, err := conn.Read(buffer)
		if err == io.EOF {
			break
		}
		response := messageHandler(string(buffer[:n]))
		conn.Write([]byte(response + "\n"))
	}
}
