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

type Message struct {
	Uuid   string
	Data string
}

type ContainerMessages struct {
    Name string
    Messages []string
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
        text := "check"
        msg := Message{Uuid: uuid.New().String(), Data: text}

        binBuf := new(bytes.Buffer)
        gobobj := gob.NewEncoder(binBuf)
        gobobj.Encode(msg)

        c.Write(binBuf.Bytes())

         tmp := make([]byte, 500)
        c.Read(tmp)

        tmpbuff := bytes.NewBuffer(tmp)
        tmpstruct := new(ContainerMessages)
        // creates a decoder object
        gobobjdec := gob.NewDecoder(tmpbuff)
        // decodes buffer and unmarshals it into a Message struct
        gobobjdec.Decode(tmpstruct)
        for _, message := range tmpstruct.Messages {
            fmt.Println(message)
        }

//         message, _ := bufio.NewReader(c).ReadString('\n')
//         if strings.TrimSpace(string(message)) == "STOP" {
//             fmt.Println("TCP server exiting...")
//             return
//         }
//
//         if strings.TrimSpace(string(message)) == "PING" {
//             continue
//         }

        //ContainerMessages in message
        //var containerMessages ContainerMessages;
        //binBuf = bytes.NewBuffer([]byte(message))
        //gobobjdec := gob.NewDecoder(binBuf)
        //gobobjdec.Decode(&containerMessages)
       // fmt.Println(tmpstruct)
        //fmt.Println(containerMessages.Name)
        //fmt.Println(containerMessages.Messages)

        time.Sleep(1 * time.Second)
    }
}
