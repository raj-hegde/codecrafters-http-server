package main

import (
	"fmt"
	"strings"

	"net"
	"os"
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	buf := make([]byte, 512)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			return
		}
		req := string(buf[:n])
		split_req := strings.Split(req, " ")
		switch {
		case split_req[1] == "/":
			conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
		case strings.HasPrefix(split_req[1], "/echo"):
			message := strings.Replace(split_req[1], "/echo/", "", 1)
			res := response(message)
			conn.Write([]byte(res))
		case split_req[1] == "/user-agent":
			message := strings.Split(split_req[len(split_req)-2], "\r\n")
			res := response(message[0])
			conn.Write([]byte(res))
		case strings.HasPrefix(split_req[1], "/files"):
			args := os.Args
			directory := args[2]
			filename := strings.Replace(split_req[1], "/files/", "", 1)
			if split_req[0] == "GET" {
				content, err := os.ReadFile(directory + "/" + filename)
				if err != nil {
					conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
					os.Exit(1)
				}
				res := response1(string(content))
				conn.Write([]byte(res))
			}
			if split_req[0] == "POST" {
				path := directory + "/" + filename
				split := strings.Split(req, "\r\n")
				content := split[len(split)-1]
				err := os.WriteFile(path, []byte(content), 0666)
				if err != nil {
					fmt.Println("error writing to file", err.Error())
					os.Exit(1)
				}
				res := response2(content)
				conn.Write([]byte(res))
			}
		default:
			conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))

		}
	}
}
func response(message string) string {
	statusLine := "HTTP/1.1 200 OK\r\n"
	contentType := "Content-Type: text/plain\r\n"
	contentLength := fmt.Sprintf("content-Length: %d\r\n\r\n", len(message))

	return statusLine + contentType + contentLength + message
}

func response1(message string) string {
	statusLine := "HTTP/1.1 200 OK\r\n"
	contentType := "Content-Type: application/octet-stream\r\n"
	contentLength := fmt.Sprintf("content-Length: %d\r\n\r\n", len(message))

	return statusLine + contentType + contentLength + message
}

func response2(message string) string {
	statusLine := "HTTP/1.1 201 OK\r\n"
	contentType := "Content-Type: application/octet-stream\r\n"
	contentLength := fmt.Sprintf("content-Length: %d\r\n\r\n", len(message))

	return statusLine + contentType + contentLength + message
}
