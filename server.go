package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
)

// All request types for http 1.0
const POST = "POST"
const GET = "GET"
const HEAD = "HEAD"

// Will parse and return a map of the request headers
func parseRequestHeader(reader *bufio.Reader) map[string][]byte {
	requestHeaders := make(map[string][]byte)
	for {
		line, err := reader.ReadBytes('\n')

		if err != nil {
			break
		}
		line = bytes.TrimSuffix(line, []byte("\r"))
		// Check for empty line to sepeate request header from body
		if len(line) == 1 && line[0] == '\r' {
			break
		}
		parts := bytes.SplitN(line, []byte(":"), 2)
		fmt.Println(parts)

		if len(parts) != 2 {
			fmt.Println("Invalid header line")
			break
		}

		key := bytes.TrimSpace(parts[0])
		value := bytes.TrimSpace(parts[1])
		requestHeaders[string(key)] = value
	}
	return requestHeaders
}

func sendResponse(conn net.Conn, response Response) {
	stringResponse := fmt.Sprintf("")
	conn.Write([]byte(stringResponse))
}

func readBody(reader *bufio.Reader, contentLength int) []byte {
	body := make([]byte, contentLength)
	_, err := io.ReadFull(reader, body)
	if err != nil {
		fmt.Println(err)
	}
	return body
}

func getContentLength(contentLength []byte) int {
	contentLengthStr := string(contentLength)
	if contentLengthStr == "" {
		return 0
	}
	length, err := strconv.Atoi(contentLengthStr)
	if err != nil {
		fmt.Println(err)
	}
	return length
}

type Request struct {
	Method  string
	Path    string
	Version string
	Headers map[string][]byte
	Body    []byte
}

type Response struct {
	Headers    map[string][]byte
	StatusCode int
	Body       []byte
}

type Router struct {
	Routes map[string]func(Request) Response
}

func InitRouter() *Router {
	return &Router{
		Routes: make(map[string]func(Request) Response),
	}
}

func (r *Router) AddRoute(method, path string, handler func(Request) Response) {
	key := method + " " + path
	r.Routes[key] = handler
}

func (r *Router) Route(req Request) Response {
	key := req.Method + " " + req.Path

	if handler, exists := r.Routes[key]; exists {
		return handler(req)
	}

	return Response{StatusCode: 404, Body: []byte("Not Found")}
}

func handleConnection(conn net.Conn, router *Router) {
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
	headerMap := parseRequestHeader(reader)
	contentLength := getContentLength(headerMap["Content-Length"])
	var body []byte

	if method == POST {
		body = readBody(reader, contentLength)
	}
	request := Request{
		Method:  method,
		Path:    path,
		Version: version,
		Headers: headerMap,
		Body:    body,
	}
	response := router.Route(request)

	sendResponse(conn, response)
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
	router := InitRouter()
	router.AddRoute("POST", "/blog", printAPI)
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		go handleConnection(conn, router)
	}
}
