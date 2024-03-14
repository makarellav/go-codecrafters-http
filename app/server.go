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

	for {
		conn, err := l.Accept()

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

		reqFields := strings.Split(string(data), "\r\n")
		startLine := reqFields[0]

		path := strings.Fields(startLine)[1]

		switch {
		case path == "/":
			conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
		case path == "/user-agent":
			userAgent, _ := strings.CutPrefix(reqFields[2], "User-Agent: ")

			conn.Write([]byte(fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s\r\n", len(userAgent), userAgent)))
		case strings.HasPrefix(path, "/echo/"):
			randomString, _ := strings.CutPrefix(path, "/echo/")

			conn.Write([]byte(fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s\r\n", len(randomString), randomString)))
		default:
			conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
		}

		conn.Close()
	}
}
