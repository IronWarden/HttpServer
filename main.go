package main

import (
	"fmt"
	"http/http"
)

func main() {
	fmt.Println("Starting Server...")
	router := http.InitRouter()
	router.AddRoute("/api", http.PrintAPI)
	address := "127.0.0.1:8080"
	http.Listen(router, address)
}
