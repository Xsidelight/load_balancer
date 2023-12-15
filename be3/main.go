package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Received request from %s\n", r.RemoteAddr)
		fmt.Fprintf(w, "Hello From Backend Server")
		fmt.Println("Replied with a hello message")
	})

	http.ListenAndServe(":8082", nil)
}
