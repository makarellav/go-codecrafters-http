package main

import (
	"fmt"
	"net"
	"os"
	"strings"
	// Uncomment this block to pass the first stage
	// "net"
	// "os"
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	l, err := net.Listen("tcp", "0.0.0.0:4221")
	defer l.Close()

	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	conn, err := l.Accept()
	defer conn.Close()

	if err != nil {
		fmt.Println("Error accepting connection ", err.Error())
		os.Exit(1)
	}

	data := make([]byte, 1024)

	_, err = conn.Read(data)

	if err != nil {
		fmt.Println("Failed to read data ", err.Error())
		os.Exit(1)
	}

	startLine := strings.Split(string(data), "\r\n")[0]

	path := strings.Fields(startLine)[1]

	if path == "/" {
		conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
	} else if strings.HasPrefix(path, "/echo/") {
		randomString, _ := strings.CutPrefix(path, "/echo/")

		conn.Write([]byte(fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s\r\n", len(randomString), randomString)))
	} else {
		conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
	}
}
