package server

import (
	"bufio"
	"net"
	"testing"
	"time"
)

func TestServer(t *testing.T) {
	expected := "Hello, John\n"

	SetMessageHandler(func(user string) string {
		return "Hello, " + user
	})
	go Listen()
	time.Sleep(10 * time.Millisecond) // Wait for server initialization

	// Connect to server and send message
	conn, err := net.Dial("tcp", ":7123")
	if err != nil {
		t.Errorf("Error connecting to server: %v", err)
	}
	conn.Write([]byte("John"))

	// Verify response
	response, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		t.Errorf("Error reading reponse from server: %v", err)
	}
	if response != expected {
		t.Errorf("Incorrect result, expected response to be: %v, got %v", expected, response)
	}
}
