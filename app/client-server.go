// Package main is intended to provide a working example of how to properly
// read from a net.Conn without losing data or integrity. The data is generated
// from alphabetical characters rather than simply rand to prevent the
// generation of an `EOF`
package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"math/rand"
	"net"
	"time"
)

var (
	done      = make(chan bool) // channel for holding exit until server finished
	clientMD5 [16]byte          // holds md5sum of data client sent
	serverMD5 [16]byte          // holds md5sum of data server received
)

func main() {
	// define the server
	serverSocket, err := net.Listen("tcp", "127.0.0.1:8082")
	if err != nil {
		fmt.Printf("Server failed to start - %s\n", err)
		return
	}
    c := make(chan []byte)
	// start the server
	go func() {
		for {
			conn, err := serverSocket.Accept()
			if err != nil {
				fmt.Printf("Server failed to accept connection - %s\n", err)
				return
			}

			// handle the connection (read the data)
			go handleConnection(conn, c)
		}
	}()

	// give the server time to start
	time.Sleep(time.Second * 2)

	// get random data
	data := randBytes(1024 * 1024)

	// store data md5
	clientMD5 = md5.Sum(data)

	// dial the server
	conn, err := net.Dial("tcp", "127.0.0.1:8082")
	if err != nil {
		fmt.Printf("Dial error - %s\n", err)
		return
	}

	// write the data to the server
	length, err := conn.Write(data)
	if err != nil {
		fmt.Printf("Write error - %s\n", err)
		return
	}

	// log the bytes written
	fmt.Printf("WROTE %d bytes\n", length)

	// wait for the connection handler to finish
	//<-done

	dataChanel := <-c
	fmt.Println(string(dataChanel))

	// log it all
	fmt.Printf("CLIENT MD5 - %x\n", clientMD5)
	fmt.Printf("SERVER MD5 - %x\n", serverMD5)
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
	}

	// store data md5
	serverMD5 = md5.Sum(data)

	// log bytes read
	fmt.Printf("READ  %d bytes\n", length)
    c <- data
	//done <- true
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func randBytes(n int) []byte {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return b
}
