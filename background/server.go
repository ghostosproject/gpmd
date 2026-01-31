package background

import (
	"log"
	"net"
	"os"

	"github.com/ghostosproject/gpmd/module"
)

func Server(mod *module.ModuleService) {
	listener, err := net.Listen("tcp", ":1111")
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
		go HandleBackgroundConnection(conn, mod)
	}
}
