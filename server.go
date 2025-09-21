package main

import (
	"fmt"
	"net"
)

func parseRequestLine(request string) {
	var requestLine string
	for i, char := range request {
		if char == '\n' {
			requestLine = request[:i]
			break
		}
	}
	fmt.Println(requestLine)
}
func parse(request []byte) {
	stringRequest := string(request)
	fmt.Println(stringRequest)
	parseRequestLine(stringRequest)
}

func main() {
	// Step 1: Establish the TCP Connection
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Listening on port 8080")
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		request := make([]byte, 1024)
		if _, err := conn.Read(request); err != nil {
			fmt.Println(err)
			continue
		} else {
			go parse(request)
		}
	}
}
