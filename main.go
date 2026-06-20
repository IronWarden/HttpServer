package main

import (
	"fmt"
	"http/http"
)

func main() {
	fmt.Println("Starting Server...")
	router := http.InitRouter()
	router.AddRoute("/api", http.PrintAPI)
	http.Listen(router)
}
