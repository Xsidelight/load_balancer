package main

import (
	"fmt"
	"io"
	"net"
)

func handleCOnnection(conn net.Conn, backendAddr string) {
	defer conn.Close()

	// Connect to backend server
	backendConn, err := net.Dial("tcp", backendAddr)
	if err != nil {
		fmt.Printf("Error connecting to backend: %s\n", err)
		return
	}

	defer backendConn.Close()

	// Forward request to backend
	go io.Copy(backendConn, conn)

	// Forward response to client
	io.Copy(conn, backendConn)

	fmt.Println("Response from server sent back to client")
}

func main() {
	lbAddr := ":80"
	backendAddr := "localhost:8080"

	listener, err := net.Listen("tcp", lbAddr)
	if err != nil {
		fmt.Printf("Error listening on %s: %s\n", lbAddr, err)
		return
	}

	defer listener.Close()

	fmt.Printf("Load Balancer listening on %s\n", lbAddr)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Error accepting connection: %s\n", err)
			return
		}

		fmt.Printf("Connection accepted from %s\n", conn.RemoteAddr())
		go handleCOnnection(conn, backendAddr)
	}

}
