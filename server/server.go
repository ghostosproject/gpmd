package server

import (
	"log"
	"net"
	"os"

	"github.com/ghostosproject/gpmd/connection"
)

func Server() {
	listener, err := net.Listen("tcp", ":5989")
	if err != nil {
		log.Printf("Error setting up tcp listener: %v\n", err)
		os.Exit(1)
	}

	nodes := map[string]connection.Node{}
	conns := map[net.Conn]string{}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v", err)
			continue // Continue listening for other connections
		}
		go connection.HandleGPMDConnection(conn, &nodes, &conns)
	}
}
