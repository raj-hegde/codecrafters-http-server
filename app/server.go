package main

import (
	"fmt"
	"strings"

	// Uncomment this block to pass the first stage
	"net"
	"os"
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

	conn, err := l.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}

	handleConn(conn)
}

func handleConn(conn net.Conn) {
	buf := make([]byte, 512)
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		return
	}
	req := string(buf[:n])
	split_req := strings.Split(req, " ")
	if split_req[1] == "/" {
		conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
	} else if strings.HasPrefix(split_req[1], "/echo") {
		message := strings.Replace(split_req[1], "/echo/", "", 1)
		res := response(message)
		conn.Write([]byte(res))
	} else if split_req[1] == "/user-agent" {
		message := strings.Split(split_req[len(split_req)-2], "\r\n")
		res := response(message[0])
		conn.Write([]byte(res))
	} else {
		conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))

	}
}

func response(message string) string {
	statusLine := "HTTP/1.1 200 OK\r\n"
	contentType := "Content-Type: text/plain\r\n"
	contentLength := fmt.Sprintf("content-Length: %d\r\n\r\n", len(message))

	return statusLine + contentType + contentLength + message
}
