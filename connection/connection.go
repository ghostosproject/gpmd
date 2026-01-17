package connection

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/ghostosproject/ghost-network-common/gmp"
)

type Connection struct {
	Type string `json:"type"`
	Node Node   `json:"node"`
}

type Node struct {
	Id      string `json:"id"`
	Name    string `json:"name"`
	Address string `json:"address"`
	Port    string `json:"port"`
	PID     int32  `json:"pid"`
	Active  bool
	conn    net.Conn
}

func HandleGPMDConnection(conn net.Conn, nodes *map[string]Node, conns *map[net.Conn]string) {
	defer conn.Close()

	addr := strings.Split(conn.RemoteAddr().String(), ":")
	fmt.Printf("Address: %s\n", addr[0])
	fmt.Printf("Port: %s\n", addr[1])

	nds := *nodes
	cns := *conns

	buffer := make([]byte, 1024) // Create a buffer of 1024 bytes

	for {
		n, err := conn.Read(buffer)
		if err != nil {
			if err.Error() == "EOF" { // Check for end of file/connection closed
				fmt.Println("Client disconnected.")
				name := cns[conn]
				node := nds[name]
				node.Active = false
				node.conn = nil
				nds[name] = node
				delete(cns, conn)
			} else {
				fmt.Printf("Error reading from connection: %v\n", err)
				name := cns[conn]
				node := nds[name]
				node.Active = false
				node.conn = nil
				nds[name] = node
				delete(cns, conn)
				fmt.Println("Client disconnected.")
			}
			return
		}

		var msg gmp.Message
		json.Unmarshal(buffer[:n], &msg)
		if msg.Type == 1 { // a GMPD connection
			// look at msg code
			if msg.Code == 1 {
				// register node and send back connection accepted
				addr := strings.Split(conn.RemoteAddr().String(), ":")

				// decode body
				var connection Connection

				err := json.Unmarshal([]byte(msg.Body), &connection)
				node := connection.Node
				if nds[node.Name].Active {
					ret := gmp.Message{
						Type: 1,
						Code: 4,
						Body: "A Node already is connected with that name..",
					}
					retMsg, err := json.Marshal(ret)
					if err != nil {
						log.Printf("Message Error: %v", err)
						continue
					}
					conn.Write(retMsg)
					return
				}
				node.Address = addr[0]
				node.Port = addr[1]
				node.conn = conn
				node.Active = true

				// check if node with the same name exists on the network

				// save nodes and connection information
				nds[node.Name] = node
				nodes = &nds

				cns[conn] = node.Name
				conns = &cns

				// if successful
				ret := gmp.Message{
					Type: 1,
					Code: 3,
				}
				retMsg, err := json.Marshal(ret)
				if err != nil {
					log.Printf("Message Error: %v", err)
					continue
				}
				conn.Write(retMsg)
			}
		}

	}
}
