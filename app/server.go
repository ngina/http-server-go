package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
	"strings"
)

func handleConnection(conn net.Conn) {
	// Create a buffer to hold incoming requests
	request := make([]byte, 1024)
	fmt.Println("Waiting for request...")

	readBytes, err := conn.Read(request)
	if err != nil {
		if err != io.EOF {
			fmt.Println("Error reading request:", err.Error())
			os.Exit(1)
		}
		fmt.Println("Received EOF: Connection closed by client. Exiting...")
		return
	}
	fmt.Println("Received request:", string(request[:readBytes]))

	// Get path and respond with 200 OK or 404 Not Found if path is "/"
	const okStatusLine = "HTTP/1.1 200 OK\r\n"
	path := strings.Split(string(request[:readBytes]), " ")[1]
	if path == "/" {
		headers := "\r\n"
		fmt.Println("Sending response:", okStatusLine+headers)
		conn.Write([]byte(okStatusLine + headers))

	} else if strings.HasPrefix(path, "/echo") {
		body := strings.Split(path[1:], "/")[1] //todo: validate path is /echo/<body>
		fmt.Println("Echoing body:", body, len(body))
		response := fmt.Sprintf(okStatusLine+"Content-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(body), body)
		conn.Write([]byte(response))

	} else if strings.HasPrefix(path, "/user-agent") {
		requestHeaders := strings.Split(string(request[:readBytes]), "\r\n")
		fmt.Println("Request strings split by \\r\\n:", requestHeaders, len(requestHeaders))

		headers := requestHeaders[1:4]
		var userAgentHeader string
		for _, header := range headers {
			if strings.HasPrefix(header, "User-Agent") { //todo: validate header is User-Agent
				userAgentHeader = header
				break
			}
		}
		body := strings.Split(userAgentHeader, ": ")[1]
		fmt.Println("Echoing body:", body)
		responseHeadersAndBody := fmt.Sprintf("Content-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(body), body)
		response := okStatusLine + responseHeadersAndBody
		conn.Write([]byte(response))

	} else if strings.HasPrefix(path, "/files") {
		filename := strings.Split(path[1:], "/")[1]
		fmt.Printf("Received file: %s\n", filename)

		file, err := os.Open(filepath.Join("/tmp/data/codecrafters.io/http-server-tester", filename))
		if err != nil {
			if os.IsNotExist(err) {
				fmt.Printf("File with name [%s] not found", filename)
				conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
			}
			log.Fatal(err)
		}
		fmt.Printf("Opened file: %s\n", filename)

		data := make([]byte, 100)
		count, err := file.Read(data)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("read %d bytes: %q\n", count, data[:count])
		responseHeadersAndBody := fmt.Sprintf("Content-Type: application/octet-stream\r\nContent-Length: %d\r\n\r\n%s", count, data[:count])
		conn.Write([]byte(okStatusLine + responseHeadersAndBody))

	} else {
		conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
	}
	fmt.Println("Sent response")
}

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	// Create listener
	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}
	defer l.Close()

	// Listen for connections
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
		fmt.Println("Connection accepted")
		go handleConnection(conn)
		fmt.Println("Connection handled")
	}
}
