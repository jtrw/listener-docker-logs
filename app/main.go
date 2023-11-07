package main

import (
	"fmt"
	"github.com/jessevdk/go-flags"
	"log"
	"os"
	"gopkg.in/yaml.v3"
	"os/exec"
	"regexp"
	"time"
	"net"
	"bytes"
	"encoding/gob"
	"strings"
)


var revision string = "1.0"

type Listener struct {
    Containers []struct {
        Name string `yaml:"name"`
        Regexp []string `yaml:"regexp"`
        Label string `yaml:"label"`
    } `yaml:"containers"`
}

type Options struct {
	Config string  `short:"f" long:"file" env:"CONF" default:"listener.yml" description:"config file"`
	Port string `short:"p" long:"port" env:"PORT" default:"2323" description:"port"`
}

type FondMessages struct {
    Container []ContainerMessages
}

type ContainerMessages struct {
    Name string
    Messages []string
}

type TcpServer struct {
    Port int
}

type Message struct {
	Uuid   string
	Data string
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

    port := ":" + opts.Port
    l, err := net.Listen("tcp", port)
    if err != nil {
        fmt.Println(err)
        return
    }
    defer l.Close()

    c, err := l.Accept()
    if err != nil {
        fmt.Println(err)
        return
    }

    for {
        tmp := make([]byte, 500)
        c.Read(tmp)
        tmpbuff := bytes.NewBuffer(tmp)
        tmpstruct := new(Message)
        // creates a decoder object
        gobobjdec := gob.NewDecoder(tmpbuff)
        // decodes buffer and unmarshals it into a Message struct
        gobobjdec.Decode(tmpstruct)
        fmt.Println(tmpstruct.Data)
        if err != nil {
            fmt.Println(err)
            return
        }
        if strings.TrimSpace(string(tmpstruct.Data)) == "STOP" {
            fmt.Println("Exiting TCP server!")
            return
        }


        for _, container := range listener.Containers {
            var containerMessages ContainerMessages;
            containerMessages.Name = string(container.Name)
           // cmd := exec.Command("docker", "logs", string(container.Name), "--tail", "30")
            time := time.Now().Add(-time.Second * 1).Format("2006-01-02T15:04:05")
            fmt.Println(time)
            cmd := exec.Command("docker", "logs", string(container.Name), "--since", time)

            output, err := cmd.CombinedOutput()

            if err != nil {
                log.Fatal(err)
            }
            outStr := string(output)

            for _, regExpStr := range container.Regexp {
                matched := regexp.MustCompile(regExpStr)
                matches := matched.FindAllStringSubmatch(outStr, -1)
                //matchesIndexes := matched.FindAllStringSubmatchIndex(outStr, -1)

                for _, v := range matches {
                    containerMessages.Messages = append(containerMessages.Messages, v[1])
                }
            }
            if len(containerMessages.Messages) > 0 {
                binBuf := new(bytes.Buffer)
                gobobj := gob.NewEncoder(binBuf)
                gobobj.Encode(containerMessages)
                fmt.Println("Send to Client");
                fmt.Println(containerMessages)
                c.Write(append(binBuf.Bytes(), '\n'))
            } else {
                c.Write([]byte("PING\n"))
            }
        }
        //time.Sleep(10 * time.Second)
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

