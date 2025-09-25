package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	// Libraries for HTTP server:
	// "net/http" (for higher-level HTTP handling, though we're implementing a basic one)
	// "io" (for input/output operations)
)

func parseRequestLine(request string) {
	// Pseudocode:
	// 1. Split the request line into method, path, and protocol version.
	// 2. Validate the method (e.g., GET, POST).
	// 3. Validate the path.
	// 4. Validate the protocol version (e.g., HTTP/1.0).
	var requestLine string
	fmt.Println(requestLine)
}

func parse(request []byte) {
	// Pseudocode:
	// 1. Convert byte slice to string.
	// 2. Iterate through the string to find the request line (first line).
	// 3. Call parseRequestLine with the extracted request line.
	// 4. Parse headers (key-value pairs) until an empty line is encountered.
	// 5. If it's a POST request, read the body based on Content-Length header.
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

// All request types for http 1.0
const POST = "POST"
const GET = "GET"
const HEAD = "HEAD"

// All content types
const html = "text/html"
const plain = "text/plain"
const gif = "image/gif"
const jpeg = "image/jpeg"
const octet = "application/octet-stream"
const form = "application/x-www-form-urlencoded"

// Will parse and return a map of the request headers
func parseRequestHeader(reader *bufio.Reader) map[string]string {
	requestHeaders := make(map[string]string)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		// Check for empty line to sepeate request header from body
		if strings.TrimSpace(line) == "" {
			break
		}
		parts := strings.Split(strings.TrimSpace(line), ":")
		if len(parts) != 2 {
			fmt.Println("Invalid header line")
			break
		}
		requestHeaders[parts[0]] = parts[1]
	}
	return requestHeaders
}

func parseRequestBody(reader *bufio.Reader) {

}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	// Pseudocode:
	// 1. Create a new buffered reader for the connection.
	// 2. Read the request line from the client.
	// 3. Handle any errors during read (e.g., connection closed).
	// 4. Parse the request line to extract method, path, and HTTP version.
	// 5. Read and parse HTTP headers.
	// 6. If necessary (e.g., POST request), read the request body.
	// 7. Determine the appropriate response based on the request (e.g., route to a handler).
	// 8. Construct the HTTP/1.0 response (Status-Line, Headers, Body).
	// 9. Write the response back to the client.
	reader := bufio.NewReader(conn)
	requestLine, err := reader.ReadString('\n')

	if err != nil {
		fmt.Println(err)
	}

	parts := strings.Split(strings.TrimSpace(requestLine), " ")
	if len(parts) != 3 {
		fmt.Println("Invalid request line")
	}
	method, path, version := parts[0], parts[1], parts[2]

	parseRequestHeader(reader)
	parseRequestBody(reader)

}
func main() {
	// Pseudocode:
	// 1. Choose a port to listen on (e.g., 8080).
	// 2. Establish a TCP listener on the chosen port.
	// 3. Handle any errors during listener creation.
	// 4. Continuously accept incoming TCP connections in a loop.
	// 5. For each accepted connection, spin up a new goroutine to handle it concurrently.
	// 6. In the goroutine, implement the logic to:
	//    a. Read the HTTP request.
	//    b. Parse the request (method, path, headers, body).
	//    c. Generate an appropriate HTTP response.
	//    d. Write the response back to the client.
	//    e. Close the connection.
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
