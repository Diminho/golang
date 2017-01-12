package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
)

const (
	CONN_TYPE = "tcp"
)

var proxyAddr *string
var remoteAddr *string
var response string

func main() {
	proxyAddr = flag.String("proxy", "", "Address [:port] of proxy server. (Required)")
	remoteAddr = flag.String("remote", "", "Address [:port] of remote server. (Required)")
	flag.Parse()
	if *proxyAddr == "" || *remoteAddr == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}
	startServer()
}

func startServer() {
	// Listen for incoming connections.
	l, err := net.Listen(CONN_TYPE, *proxyAddr)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	// Close the listener when the application closes.
	defer l.Close()
	fmt.Println("Listening on " + *proxyAddr)
	for {
		// Listen for an incoming connection.
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}
		// Handle connections in a new goroutine.
		go handleRequest(conn)
	}
}

// Handles incoming requests.
func handleRequest(conn net.Conn) {
	// Make a buffer to hold incoming data.
	buf := make([]byte, 1024)

	// Read the incoming connection into the buffer.
	length, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
	}

	remoteConnection, remoteErr := net.Dial(CONN_TYPE, *remoteAddr)
	if remoteErr != nil {
		fmt.Println(remoteErr)
		os.Exit(1)
	}
	remoteConnection.Write([]byte(string(buf[:length])))
	connbuf := bufio.NewReader(remoteConnection)
	for {
		response, err := connbuf.ReadString('\n')
		if len(response) > 0 {
			fmt.Println(response)
		}
		if err != nil {
			break
		}
	}
	// Send a response back to person contacting us.
	conn.Write([]byte(string(response)))
	// Close the connection when you're done with it.
	conn.Close()
}
