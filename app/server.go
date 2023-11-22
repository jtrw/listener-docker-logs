package main

import (
	"encoding/gob"
	"fmt"
	"net"
)

type ContainerMessages struct {
	Name     string
	Messages []string
}

type Message struct {
	Uuid    string
	Data    ContainerMessages
	Message string
}

func handleClientConnection(conn net.Conn) {
	defer conn.Close()

	// Create a new encoder to send gob data
	encoder := gob.NewEncoder(conn)

	// Example data to send
	dataToSend := Message{
		Uuid: "123",
		Data: ContainerMessages{
			Name:     "Example",
			Messages: []string{"Message 1", "Message 2", "Message 3"},
		},
		Message: "Hello from server!",
	}

	// Send data to client
	err := encoder.Encode(dataToSend)
	if err != nil {
		fmt.Println("Error encoding and sending data:", err)
		return
	}
}

func main() {
	// Start server
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	defer listener.Close()
	fmt.Println("Server started. Waiting for clients...")

	// Accept client connections
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		fmt.Println("Client connected:", conn.RemoteAddr())

		// Handle client connection in a goroutine
		go handleClientConnection(conn)
	}
}
