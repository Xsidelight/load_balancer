package main

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

var (
	backendServers = []string{"127.0.0.1:8080", "127.0.0.1:8081", "127.0.0.1:8082"}
	currentServer  uint64
	healthChecks   sync.Once
)

func getNextBackendServer() string {
	serverIndex := atomic.AddUint64(&currentServer, 1) % uint64(len(backendServers))
	return backendServers[serverIndex]
}

func healthCheck(backendAddr string) {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		resp, err := http.Get("http://" + backendAddr)
		if err != nil {
			fmt.Printf("Server %s is unavailable: %s\n", backendAddr, err)
			continue
		}
		fmt.Printf("Server %s is healthy\n", backendAddr)
		resp.Body.Close()
	}
}

func startHealthChecks() {
	for _, addr := range backendServers {
		go healthCheck(addr)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	backendAddr := getNextBackendServer()

	backendConn, err := net.Dial("tcp", backendAddr)
	if err != nil {
		fmt.Printf("Error connecting to backend: %s\n", err)
		return
	}
	defer backendConn.Close()

	// Use io.Copy for bi-directional copying
	if _, err := io.Copy(backendConn, conn); err != nil {
		fmt.Println("Error forwarding request:", err)
		return
	}
	if _, err := io.Copy(conn, backendConn); err != nil {
		fmt.Println("Error sending response:", err)
	}
}

func main() {
	// Start health checks once for all servers
	healthChecks.Do(startHealthChecks)

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
			continue
		}

		go handleConnection(conn)
	}
}
