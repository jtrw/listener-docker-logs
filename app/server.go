package main

import (
	"encoding/gob"
	"fmt"
	"net"
	"gopkg.in/yaml.v3"
	"github.com/jessevdk/go-flags"
	"log"
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

var revision string = "1.0"

type Listener struct {
    Containers []struct {
        Name string `yaml:"name"`
        Regexp []string `yaml:"regexp"`
        Label string `yaml:"label"`
        All bool `yaml:"all"`
    } `yaml:"containers"`
}

type Options struct {
	Config string  `short:"f" long:"file" env:"CONF" default:"listener.yml" description:"config file"`
	Port string `short:"p" long:"port" env:"PORT" default:"2323" description:"port"`
}

func main() {
    fmt.Printf("Listener %s\n", revision)

    var opts Options
    parser := flags.NewParser(&opts, flags.Default)
    _, err := parser.Parse()
    if err != nil {
        log.Fatal(err)
    }

    listener, errYaml := LoadConfig(opts.Config)
    if errYaml != nil {
        log.Println(errYaml)
    }

	// Start server
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	defer listener.Close()
	fmt.Println("Server started. Waiting for clients...")

    tmpTime := time.Now().Add(-time.Second * 10).Format("2006-01-02T15:04:05")

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
