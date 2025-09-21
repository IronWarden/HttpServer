package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

func parseRequestLine(request string) {
	var requestLine string
	fmt.Println(requestLine)
}

func parse(request []byte) {
	stringRequest := string(request)
	fmt.Println(stringRequest)
	isRequestLine := false

	for i, char := range stringRequest {
		if char == '\n' && !isRequestLine {
			requestLine := stringRequest[:i]
			parseRequestLine(requestLine)
			isRequestLine = true
			stringRequest = stringRequest[i+1:]
		}

	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	requestLine, err := reader.ReadString('\n')

	if err != nil {
		fmt.Println(err)
	}

	parts := strings.Split(strings.TrimSpace(requestLine), " ")
	if len(parts) != 3 {
		fmt.Println("Invalid request line")
	}

}
func main() {
	// Step 1: Establish the TCP Connection
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println(err)
	}
	defer listener.Close()
	fmt.Println("Listening on port 8080")
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		go handleConnection(conn)
	}
}
