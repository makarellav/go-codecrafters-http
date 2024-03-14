package main

import (
	"bytes"
	"fmt"
	"net"
	"os"
	// Uncomment this block to pass the first stage
	// "net"
	// "os"
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	l, err := net.Listen("tcp", "0.0.0.0:4221")

	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	conn, err := l.Accept()

	if err != nil {
		fmt.Println("Error accepting connection ", err.Error())
		os.Exit(1)
	}

	var data []byte

	_, err = conn.Read(data)

	if err != nil {
		fmt.Println("Failed to read data ", err.Error())
		os.Exit(1)
	}

	_, err = conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))

	if err != nil {
		fmt.Println("Failed to write data", err.Error())
		os.Exit(1)
	}

	fmt.Println(string(data))

	startLine := bytes.Split(data, []byte("\r\n"))[0]
	path := bytes.Split(startLine, []byte(" "))[1]

	fmt.Println(string(startLine))
	fmt.Println(string(path))
}
