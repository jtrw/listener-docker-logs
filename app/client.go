package main

import (
	"encoding/gob"
	"fmt"
	"net"
	"time"
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

func main() {
	for {
		// Connect to server
		conn, err := net.Dial("tcp", "127.0.0.1:8080")
		if err != nil {
			fmt.Println("Error connecting to server:", err)
			time.Sleep(1 * time.Second) // Почекати перед спробою нового з'єднання
			continue
		}

		// Create a new decoder to receive gob data
		decoder := gob.NewDecoder(conn)

		// Receive data from server
		var receivedData Message
		err = decoder.Decode(&receivedData)
		if err != nil {
			fmt.Println("Error receiving and decoding data:", err)
			conn.Close()
			time.Sleep(1 * time.Second) // Почекати перед спробою нового з'єднання
			continue
		}

		// Print received data
		fmt.Println("Received data from server:")
		fmt.Println("UUID:", receivedData.Uuid)
		fmt.Println("Message:", receivedData.Message)
		fmt.Println("Data Name:", receivedData.Data.Name)
		fmt.Println("Data Messages:", receivedData.Data.Messages)

		conn.Close()

		time.Sleep(1 * time.Second) // Почекати 1 секунду перед наступним запитом
	}
}
