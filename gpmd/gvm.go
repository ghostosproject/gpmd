package gpmd

import (
	"fmt"
	"log"
	"net"
	"os"
)

func runGVMServer(md *GPMD, port string) {
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
		// }
		go handleGVMConnection(md, conn)
	}
}
