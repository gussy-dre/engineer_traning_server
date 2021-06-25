package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var (
	portOption = flag.Int("p", 8089, "help message for \"p\" option")
)

func handleResolveAddress(port string) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", "0.0.0.0:"+port)
	if err != nil {
		log.Println("ResolveTCPAddr", err)
		log.Println("handleResolveAddress", err)
		return
	}

	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		log.Println("ListenTCP", err)
		log.Println("handleResolveAddress", err)
		return
	}

	err = handleListener(listener)
	if err != nil {
		log.Println("handleListener", err)
		log.Println("handleResolveAddress", err)
		return
	}
}

func handleListener(listener *net.TCPListener) error {
	defer listener.Close()

	for {
		connection, err := listener.AcceptTCP()
		if err != nil {
			if ne, ok := err.(net.Error); ok {
				if ne.Temporary() {
					log.Println("AcceptTCP", err)
					continue
				}
			}
			log.Println("Unrecoverable", err)
			return err
		}
		go handleConnection(connection)
	}
}

func handleConnection(connection *net.TCPConn) {
	defer connection.Close()

	fmt.Println(">>> start")
	createHeader(connection)
	fmt.Println(">>> end" + "\n")
}

func createHeader(connection *net.TCPConn) {
	buf := make([]byte, 4096)
	header := make(map[string]string)

	n, err := connection.Read(buf)

	if err != nil {
		log.Println("Read", err)
		return
	}

	receive := strings.TrimSpace(string(buf[:n]))
	fmt.Println(receive + "\n")

	h := strings.Split(receive, "\r\n")
	request := strings.Split(h[0], " ")
	header["method"] = request[0]
	header["path"] = request[1]

	for i, v := range h {
		if i == 0 {
			continue
		}
		headerFields := strings.Split(v, ": ")
		header[headerFields[0]] = headerFields[1]
	}
	createResponce(connection, header)
}

func createResponce(connection *net.TCPConn, header map[string]string) {
	method, ok := header["method"]
	if !ok {
		log.Println("no method found")
		return
	}

	var resp []byte
	if method == "GET" {
		path, ok := header["path"]
		if !ok {
			log.Println("no path found")
			return
		} else if path == "/favicon.ico" {
			log.Println("not create favicon.ico")
			return
		} else if path == "/" {
			path = "/sample.html"
		}

		cwd, err := os.Getwd()
		if err != nil {
			log.Println("Getwd", err)
			return
		}

		p := filepath.Join(cwd, filepath.Clean(path))
		resp, err = ioutil.ReadFile(p)
		if err != nil {
			log.Println("ReadFile", err)
			return
		}
	}

	connection.Write([]byte("HTTP/1.1 200 OK\r\n"))
	connection.Write([]byte("Content-Type: text/html\r\n"))
	connection.Write([]byte("\r\n"))
	connection.Write(resp)
}

func main() {
	flag.Parse()
	port := strconv.Itoa(*portOption)

	log.Println("Server running at 0.0.0.0:" + port + "\n")

	handleResolveAddress(port)
}
