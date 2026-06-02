package gpmd

import (
	"fmt"
	"log"
	"net"
	"os"
)

func runBackgroundServer(md *GPMD, port string) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Printf("Error setting up tcp listener: %v\n", err)
		os.Exit(1)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v", err)
			continue // Continue listening for other connections
		}
		go handleBackgroundConnection(conn, md)
	}
}

func handleBackgroundConnection(conn net.Conn, md *GPMD) {
	defer conn.Close()
	buffer := make([]byte, 1024) // Create a buffer of 1024 bytes

	for {
		n, err := conn.Read(buffer)
		if err != nil {
			if err.Error() == "EOF" { // Check for end of file/connection closed
				fmt.Println("Client disconnected.")
				return
			}
		}
		fmt.Println(string(buffer[:n]))
		conn.Close()
		return
	}
}
