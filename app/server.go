package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	// Uncomment this block to pass the first stage

	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}
	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			if err != io.EOF {
				fmt.Println("Error accepting connection: ", err.Error())
				os.Exit(1)
			}
			fmt.Println("Received EOF: Connection closed by client. Exiting...")
			break
		}
		defer conn.Close()
		fmt.Println("Connection accepted")

		// Read first request and respond with 200 OK
		request := make([]byte, 1024)
		readBytes, err := conn.Read(request)
		if err != nil {
			fmt.Println("Error reading request: ", err.Error())
			os.Exit(1)
		}
		fmt.Println("Received request:\n", string(request[:readBytes]))
		fmt.Printf("Split request using empty spaces: %s\n", strings.Split(string(request[:readBytes]), " "))

		// Get path and respond with 200 OK or 404 Not Found if path is "/"
		path := strings.Split(string(request[:readBytes]), " ")[1]
		if path == "/" {
			conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
		} else {
			conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
		}
	}
}
