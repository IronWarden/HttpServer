package main

import (
	"fmt"
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

func printAPI(r Request) Response {
	response := Response{
		StatusCode: OK,
	}
	switch r.Method {
	case "GET":
		fmt.Printf("You are fetching path : %s", r.Path)
	case "POST":
		fmt.Printf("You are posting to path: %s", r.Path)
	}
	return response
}
