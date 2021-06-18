package main

import (
        "fmt"
        "net"
        "io"
        "log"
)

func main() {
	listener, error := net.Listen("tcp", "0.0.0.0:8089")

        if error != nil {
                log.Printf("Listen", error)
                return
        }

        fmt.Printf("Server running at 0.0.0.0:8089\n\n")
        
        defer connection.Close();

	for {
		connection, error := listener.Accept()

                if error != nil {
                        log.Printf("Accept", error)
                        continue
                }

                go func(connection net.Conn) {

                        for {
                                buf := make([]byte, 4096)
                                n, error := connection.Read(buf)

                                if (error != nil) {
                                        if error == io.EOF {
                                                log.Printf("io.EOF", error)
                                                continue
                                        } else {
                                                log.Printf("Read", error)
                                                continue
                                        }
                                }

                                fmt.Printf("Client> %s \n", string(buf[:n]))

                                n, error = connection.Write([]byte(string(buf[:n])))

                                if error != nil {
                                        log.Printf("Write", error)
                                        continue
                                }

                        }
                }(connection)
	}
}