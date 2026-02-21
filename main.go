package main

import (
	"./http"
	"fmt"
)

func main() {
	fmt.Println("Starting Server...")
	router := http.InitRouter()
	router.AddRoute("/api", http.PrintAPI)
	http.Listen(router)
}
