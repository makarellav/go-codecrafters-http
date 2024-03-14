package main

import (
	"bytes"
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
	// Uncomment this block to pass the first stage
	// "net"
	// "os"
)

func main() {
	dir := flag.String("directory", ".", "Directory")

	flag.Parse()

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

		go func(c net.Conn) {
			data := make([]byte, 1024)
			defer c.Close()

			_, err = conn.Read(data)

			if err != nil {
				fmt.Println("Failed to read data ", err.Error())
				os.Exit(1)
			}

			reqFields := strings.Split(string(data), "\r\n")
			startLineParts := strings.Fields(reqFields[0])

			path := startLineParts[1]

			switch {
			case path == "/":
				conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
			case path == "/user-agent":
				userAgent, _ := strings.CutPrefix(reqFields[2], "User-Agent: ")

				conn.Write([]byte(fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s\r\n", len(userAgent), userAgent)))
			case strings.HasPrefix(path, "/echo/"):
				randomString, _ := strings.CutPrefix(path, "/echo/")

				conn.Write([]byte(fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s\r\n", len(randomString), randomString)))
			case strings.HasPrefix(path, "/files/"):
				method := startLineParts[0]

				filename, _ := strings.CutPrefix(path, "/files/")
				filepath := strings.Join([]string{*dir, filename}, string(os.PathSeparator))

				switch method {
				case "GET":
					if dir == nil {
						conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))

						return
					}

					if _, err := os.Stat(filepath); os.IsNotExist(err) {
						conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))

						return
					}

					data, err := os.ReadFile(filepath)

					if err != nil {
						fmt.Println("Failed to open a file ", filepath, err.Error())

						return
					}

					conn.Write([]byte(fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: application/octet-stream\r\nContent-Length: %d\r\n\r\n%s\r\n", len(data), string(data))))
				case "POST":
					reqBody, _ := strings.CutPrefix(reqFields[len(reqFields)-1], "\r\n")
					fileData := bytes.Trim([]byte(reqBody), "\x00")

					err = os.WriteFile(filepath, fileData, 0644)

					if err != nil {
						fmt.Println("Failed to write a file ", filepath, err.Error())
					}

					conn.Write([]byte("HTTP/1.1 201 Created\r\n"))
				}
			default:
				conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
			}
		}(conn)
	}
}
