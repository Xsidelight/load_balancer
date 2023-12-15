package main

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"sync/atomic"
)

var (
	backendServers = []string{"127.0.0.1:8080", "127.0.0.1:8081", "127.0.0.1:8082"}
	currentServer  uint64
)

func getNextBackendServer() string {
	serverIndex := atomic.AddUint64(&currentServer, 1) % uint64(len(backendServers))
	return backendServers[serverIndex]
}

func healthCheck(backendAddr string) {
	response, err := http.Get("http://" + backendAddr)
	if err != nil || response.StatusCode != 200 {
		fmt.Printf("Server is unavailable: %s\n", err)
		return
	}

}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	backendAddr := getNextBackendServer()

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
		go handleConnection(conn)
	}

}
