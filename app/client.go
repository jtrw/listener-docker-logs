package main

import (
    "bufio"
    "fmt"
    "net"
    "os"
    "strings"
    "encoding/gob"
    "bytes"
    "github.com/google/uuid"
)

type Message struct {
	Uuid   string
	Data string
}

func main() {
    arguments := os.Args
    if len(arguments) == 1 {
        fmt.Println("Please provide host:port.")
        return
    }

    CONNECT := arguments[1]
    c, err := net.Dial("tcp", CONNECT)
    if err != nil {
        fmt.Println(err)
        return
    }

    for {
        //reader := bufio.NewReader(os.Stdin)
        //fmt.Print(">> ")
        //text, _ := reader.ReadString('\n')
        text := "check"
        msg := Message{Uuid: uuid.New().String(), Data: text}

        binBuf := new(bytes.Buffer)
        gobobj := gob.NewEncoder(binBuf)
        gobobj.Encode(msg)

        c.Write(binBuf.Bytes())

        message, _ := bufio.NewReader(c).ReadString('\n')
        fmt.Print("->: " + message)
        if strings.TrimSpace(string(text)) == "STOP" {
            fmt.Println("TCP client exiting...")
            return
        }
    }
}
