package main

import (
	"fmt"
	"log"
	"net"
)

func main() {
	tcpAddr, error := net.ResolveTCPAddr("tcp", "0.0.0.0:8089")
	if error != nil {
		log.Println("ResolveTCPAddr", error)
		return
	}

	listener, error := net.ListenTCP("tcp", tcpAddr)

	if error != nil {
		log.Println("ListenTCP", error)
		return
	}

	fmt.Printf("Server running at 0.0.0.0:8089\n\n")

	defer listener.Close()

	for {
		connection, error := listener.AcceptTCP()

		if error != nil {
			log.Println("AcceptTCP", error)
			return
		}

		go func(connection net.Conn) {
			defer connection.Close()

			buf := make([]byte, 4096)
			connection.Read(buf)

			connection.Write([]byte("HTTP/1.1 200 OK\r\n"))
			connection.Write([]byte("Content-Type: text/plain\r\n"))
			connection.Write([]byte("\r\n"))
			connection.Write([]byte("Hello, world"))
		}(connection)
	}
}
