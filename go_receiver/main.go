package main

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	// host := os.Getenv("HOST")
	port := os.Getenv("PORT")
	if port == "" {
        port = "8089"
    }
	// certPath := os.Getenv("CERTPATH")
	// keyPath := os.Getenv("KEYPATH")

	address := fmt.Sprintf(":%s", port)

	// listener, err := tlsServer(address, certPath, keyPath)
	listener, err := tcpServer(address)
	if err != nil {
		fmt.Printf("err starting server, %v\n", err)
		os.Exit(1)
	}

	defer listener.Close()

	for {
		// Accept incoming connections
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Error accepting connection: %v\n", err)
			continue
		}

		// Handle each connection in a separate goroutine
		go handleConnection(conn)
	}
}

func tcpServer(address string) (net.Listener, error) {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		fmt.Printf("err starting tcp server %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("tcp Server listening on %s\n", address)

	return listener, err
}

func tlsServer(address string, certPath string, keyPath string) (net.Listener, error) {
	cert, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		return nil, err
	}

	config := &tls.Config{Certificates: []tls.Certificate{cert}}

	listener, err := tls.Listen("tcp", address, config)
	if err != nil {
		return nil, err
	}

	fmt.Printf("tls Server listening on %s\n", address)

	return listener, nil
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	clientAddr := conn.RemoteAddr().String()
	fmt.Printf("New client connected: %s\n", clientAddr)

	// Create a scanner to read messages line by line
	scanner := bufio.NewScanner(conn)

	for scanner.Scan() {
		message := strings.TrimSpace(scanner.Text())

		if message == "" {
			fmt.Println("Received from nothing")
			continue
		}

		fmt.Printf("Received from %s: %s\n", clientAddr, message)
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading from %s: %v\n", clientAddr, err)
	}

	fmt.Printf("Connection closed: %s\n", clientAddr)
}
