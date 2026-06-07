package http

import (
	"fmt"
)

func PrintAPI(r Request) Response {
	response := Response{
		StatusCode: OK,
		Body:       make([]byte, 0),
		Headers:    make(map[string][]byte),
	}
	var body string
	switch r.Method {
	case "GET":
		body = fmt.Sprintf("You are fetching path: %s", r.Path)
	case "POST":
		body = fmt.Sprintf("You are posting to path: %s", r.Path)
	}
	response.Body = []byte(body)
	return response
}
