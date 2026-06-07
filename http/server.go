package http

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
const (
	POST = "POST"
	GET  = "GET"
	HEAD = "HEAD"
)

// Status Codes for HTTP 1.0
const (
	OK                    = 200
	Created               = 201
	Accepted              = 202
	No_Content            = 204
	Moved_Permanently     = 301
	Moved_Temporarily     = 302
	Not_Modified          = 304
	Bad_Request           = 400
	Unauthorized          = 401
	Forbidden             = 403
	Not_Found             = 404
	Internal_Server_Error = 500
	Not_Implemented       = 501
	Bad_Gateway           = 502
	Service_Unavailable   = 503
)

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
		if len(line) == 1 && line[0] == '\n' {
			break
		}
		parts := bytes.SplitN(line, []byte(":"), 2)

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

// Created               = 201
// Accepted              = 202
// No_Content            = 204
// Moved_Permanently     = 301
// Moved_Temporarily     = 302
// Not_Modified          = 304
// Bad_Request           = 400
// Unauthorized          = 401
// Forbidden             = 403
// Not_Found             = 404
// Internal_Server_Error = 500
// Not_Implemented       = 501
// Bad_Gateway           = 502
// Service_Unavailable   = 503
func sendResponse(conn net.Conn, response Response) {
	switch response.StatusCode {
	case OK:
		conn.Write([]byte(fmt.Sprintf("HTTP/1.0 %d OK\r\n", response.StatusCode)))
	case Created:
		conn.Write([]byte(fmt.Sprintf("HTTP/1.0 %d Created\r\n", response.StatusCode)))
	case Accepted:
		conn.Write([]byte(fmt.Sprintf("HTTP/1.0 %d Accepted\r\n", response.StatusCode)))
	case No_Content:
		conn.Write([]byte(fmt.Sprintf("HTTP/1.0 %d No Content\r\n", response.StatusCode)))
	case Moved_Permanently:
		conn.Write([]byte(fmt.Sprintf("HTTP/1.0 %d Moved Permanently\r\n", response.StatusCode)))
	case Moved_Temporarily:
		conn.Write([]byte(fmt.Sprintf("HTTP/1.0 %d Moved Temporarily\r\n", response.StatusCode)))
	case Not_Modified:
		conn.Write([]byte(fmt.Sprintf("HTTP/1.0 %d Not Modified\r\n", response.StatusCode)))
	case Bad_Request:
		conn.Write([]byte(fmt.Sprintf("HTTP/1.0 %d Bad Request\r\n", response.StatusCode)))
	case Unauthorized:
		conn.Write([]byte(fmt.Sprintf("HTTP/1.0 %d Unauthorized\r\n", response.StatusCode)))
	case Forbidden:
		conn.Write([]byte(fmt.Sprintf("HTTP/1.0 %d Forbidden\r\n", response.StatusCode)))
	case Not_Found:
		conn.Write([]byte(fmt.Sprintf("HTTP/1.0 %d Not Found\r\n", response.StatusCode)))
	case Internal_Server_Error:
		conn.Write([]byte(fmt.Sprintf("HTTP/1.0 %d Internal Server Error\r\n", response.StatusCode)))
	case Not_Implemented:
		conn.Write([]byte(fmt.Sprintf("HTTP/1.0 %d Not Implemented\r\n", response.StatusCode)))
	case Bad_Gateway:
		conn.Write([]byte(fmt.Sprintf("HTTP/1.0 %d Bad Gateway\r\n", response.StatusCode)))
	case Service_Unavailable:
		conn.Write([]byte(fmt.Sprintf("HTTP/1.0 %d Service Unavailable\r\n", response.StatusCode)))
	}

	// Default headers
	conn.Write([]byte("Server: My Server\r\n"))
	conn.Write([]byte("Content-Length: 0\r\n"))

	for key, value := range response.Headers {
		conn.Write([]byte(fmt.Sprintf("%s: %s\r\n", key, value)))
	}
	// Blank line
	conn.Write([]byte("\r\n"))
	conn.Write(response.Body)
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

func InitRouter() *Router {
	return &Router{
		Routes: make(map[string]func(Request) Response),
	}
}

func (r *Router) AddRoute(path string, handler func(Request) Response) {
	key := path
	r.Routes[key] = handler
}

func (r *Router) Route(req Request) Response {
	key := req.Path

	if handler, exists := r.Routes[key]; exists {
		return handler(req)
	}

	return Response{StatusCode: Not_Found, Body: []byte("Not Found")}
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
		conn.Close()
		return
	}

	parts := strings.Split(strings.TrimSpace(requestLine), " ")

	if len(parts) != 3 {
		fmt.Println("Invalid request line")
		conn.Close()
		return
	}
	method, path, version := parts[0], parts[1], parts[2]
	headerMap := parseRequestHeader(reader)
	var body []byte

	if method == POST {
		contentLength := getContentLength(headerMap["Content-Length"])
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

func Listen(router *Router) {
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
		go handleConnection(conn, router)
	}
}
