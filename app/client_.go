package main

import (
    //"bufio"
    "fmt"
    "net"
    "os"
    "io"
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

    chainData := make(chan []byte)

    for {
        go handleConnection(c, chainData)

        dataChanel := <-chainData
        fmt.Println(string(dataChanel))

        time.Sleep(1 * time.Second)
    }





//     for {
//         msg := Message{Uuid: uuid.New().String(), Message: "PING"}
//
//         binBuf := new(bytes.Buffer)
//         gobobj := gob.NewEncoder(binBuf)
//         gobobj.Encode(msg)
//
//         c.Write(binBuf.Bytes())
//
//         tmp := make([]byte, 10000) // 10000 bytes
//         c.Read(tmp)
//
//         tmpbuff := bytes.NewBuffer(tmp)
//         tmpstruct := new(Message)
//         // creates a decoder object
//         gobobjdec := gob.NewDecoder(tmpbuff)
//         // decodes buffer and unmarshals it into a Message struct
//         gobobjdec.Decode(tmpstruct)
//
//         //count tmpstruct.Data
//         if len(tmpstruct.Data.Messages) > 0 {
//             fmt.Println(tmpstruct.Data)
//         }
//
//         time.Sleep(1 * time.Second)
//     }
}

func handleConnection(conn net.Conn, c chan []byte) {
	// make a temporary bytes var to read from the connection
	tmp := make([]byte, 1024)
	// make 0 length data bytes (since we'll be appending)
	data := make([]byte, 0)
	// keep track of full length read
	length := 0

	// loop through the connection stream, appending tmp to data
	for {
	    fmt.Println("PING")

	    msg := Message{Uuid: uuid.New().String(), Message: "PING"}

        binBuf := new(bytes.Buffer)
        gobobj := gob.NewEncoder(binBuf)
        gobobj.Encode(msg)

        conn.Write(binBuf.Bytes())


		// read to the tmp var
		n, err := conn.Read(tmp)
		if err != nil {
			// log if not normal error
			if err != io.EOF {
				fmt.Printf("Read error - %s\n", err)
			}
			break
		}

		// append read data to full data
		data = append(data, tmp[:n]...)

		// update total read var
		length += n

		fmt.Println(length)
	}

	// log bytes read
	fmt.Printf("READ  %d bytes\n", length)
    c <- data
	//done <- true
}
