package main

import (
	"encoding/gob"
	"fmt"
	"net"
	"os"
	"gopkg.in/yaml.v3"
	"github.com/jessevdk/go-flags"
	"os/exec"
    "regexp"
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

    listenerConfig, errYaml := LoadConfig(opts.Config)
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
		go handleClientConnection(conn, listenerConfig, tmpTime)
	}
}

func handleClientConnection(conn net.Conn, listener *Listener, tmpTime string) {
	defer conn.Close()


	// Create a new encoder to send gob data
	encoder := gob.NewEncoder(conn)

	// Example data to send
// 	dataToSend := Message{
// 		Uuid: "123",
// 		Data: ContainerMessages{
// 			Name:     "Example",
// 			Messages: []string{"Message 1", "Message 2", "Message 3"},
// 		},
// 		Message: "Hello from server!",
// 	}
    msg := Message{Uuid: "1"}
	for _, container := range listener.Containers {
        var containerMessages ContainerMessages;
        containerMessages.Name = string(container.Name)
       // cmd := exec.Command("docker", "logs", string(container.Name), "--tail", "30")
        //time := time.Now().Add(-time.Second * 1).Format("2006-01-02T15:04:05")

        cmd := exec.Command("docker", "logs", string(container.Name), "--since", tmpTime)
        tmpTime = time.Now().Format("2006-01-02T15:04:05")
        output, err := cmd.CombinedOutput()

        if err != nil {
            log.Fatal(err)
        }
        outStr := string(output)
        if container.All {
            containerMessages.Messages = append(containerMessages.Messages, outStr)
        } else {
            for _, regExpStr := range container.Regexp {
                matched := regexp.MustCompile(regExpStr)
                matches := matched.FindAllStringSubmatch(outStr, -1)
                //matchesIndexes := matched.FindAllStringSubmatchIndex(outStr, -1)

                for _, v := range matches {
                    containerMessages.Messages = append(containerMessages.Messages, v[1])
                }
            }
        }
        if len(containerMessages.Messages) > 0 {
            msg = Message{Uuid: "1", Data: containerMessages}
        } else {
            msg = Message{Uuid: "1", Message: "PONG"}
        }
    }


	// Send data to client
	err := encoder.Encode(msg)
	if err != nil {
		fmt.Println("Error encoding and sending data:", err)
		return
	}
}

func LoadConfig(file string) (*Listener, error) {
	fh, err := os.Open(file)
	if err != nil {
		return nil, fmt.Errorf("can't load config file %s: %w", file, err)
	}
	defer fh.Close()

	res := Listener{}
	if err := yaml.NewDecoder(fh).Decode(&res); err != nil {
		return nil, fmt.Errorf("can't parse config: %w", err)
	}
	return &res, nil
}

