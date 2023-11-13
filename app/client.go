package main

import (
    //"bufio"
    "fmt"
    "net"
    "os"
    //"strings"
    "encoding/gob"
    "bytes"
    "github.com/google/uuid"
    "time"
)

type ContainerMessages struct {
    Name string
    Messages []string
}

type Message struct {
	Uuid   string
	Data ContainerMessages
	Message string
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
        msg := Message{Uuid: uuid.New().String(), Message: "PING"}

        binBuf := new(bytes.Buffer)
        gobobj := gob.NewEncoder(binBuf)
        gobobj.Encode(msg)

        c.Write(binBuf.Bytes())

         tmp := make([]byte, 500)
        c.Read(tmp)

        tmpbuff := bytes.NewBuffer(tmp)
        tmpstruct := new(Message)
        // creates a decoder object
        gobobjdec := gob.NewDecoder(tmpbuff)
        // decodes buffer and unmarshals it into a Message struct
        gobobjdec.Decode(tmpstruct)

        fmt.Println(tmpstruct.Data)

        time.Sleep(1 * time.Second)
    }
}
