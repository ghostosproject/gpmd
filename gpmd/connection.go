package gpmd

import (
	"bufio"
	"fmt"
	"net"

	"github.com/libp2p/go-libp2p/core/network"
)

func handleGPMDStream(s network.Stream, md *GPMD) {
	rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))
	for {
		bts, err := rw.ReadBytes(0xFE)
		if err != nil {
			panic(err)
		}
		fmt.Println(bts)
	}
}

func handleGVMConnection(md *GPMD, conn net.Conn) {
	defer conn.Close()

	// addr := strings.Split(conn.RemoteAddr().String(), ":")
	// fmt.Printf("Address: %s\n", addr[0])
	// fmt.Printf("Port: %s\n", addr[1])

	// nds := *nodes
	// cns := *conns

	buffer := make([]byte, 1024) // Create a buffer of 1024 bytes

	for {
		n, err := conn.Read(buffer)
		if err != nil {
			if err.Error() == "EOF" { // Check for end of file/connection closed
				fmt.Println("Client disconnected.")
			} else {
				fmt.Printf("Error reading from connection: %v\n", err)
				fmt.Println("Client disconnected.")
			}
			return
		}
		fmt.Println(string(buffer[:n]))
	}
}
