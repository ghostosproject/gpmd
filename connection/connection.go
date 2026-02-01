package connection

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/ghostosproject/ghost-network-common/gmp"
	"github.com/ghostosproject/gpmd/module"
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

func HandleGPMDConnection(conn net.Conn, nodes *map[string]Node, conns *map[net.Conn]string, mod *module.ModuleService) {
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
		if msg.Type == gmp.MESSAGE_TYPE_GPMD { // a GMPD connection
			// look at msg code
			if msg.Code == gmp.GMPD_REGISTER {
				// register node and send back connection accepted
				addr := strings.Split(conn.RemoteAddr().String(), ":")

				// decode body
				var connection Connection

				err := json.Unmarshal([]byte(msg.Body), &connection)
				node := connection.Node
				if nds[node.Name].Active {
					ret := gmp.Message{
						Type: gmp.MESSAGE_TYPE_GPMD,
						Code: gmp.GMPD_REGISTER_DENIED,
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
					Type: gmp.MESSAGE_TYPE_GPMD,
					Code: gmp.GMPD_REGISTER_ACCEPT,
				}
				retMsg, err := json.Marshal(ret)
				if err != nil {
					log.Printf("Message Error: %v", err)
					continue
				}
				conn.Write(retMsg)
			} else if msg.Code == gmp.GMPD_REQUEST_NODE_CONN {
				var connection Connection

				err := json.Unmarshal([]byte(msg.Body), &connection)
				if err != nil {
					fmt.Printf("Error (connection.go:122): %v\n", err)
				}
				node := connection.Node
				fmt.Println(node.Name)
				// look for node
				respNode := nds[node.Name]
				fmt.Println(respNode)
				if respNode == (Node{}) {
					respMsg := gmp.Message{
						Type: gmp.MESSAGE_TYPE_GPMD,
						Code: gmp.GMPD_NODE_CONN_DENIED,
						Body: "Node not found",
					}
					respBts, err := json.Marshal(respMsg)
					if err != nil {
						fmt.Println("Error (connection.go:136): ", err)
					}
					conn.Write(respBts)
					continue
				}
				if !respNode.Active {
					respMsg := gmp.Message{
						Type: gmp.MESSAGE_TYPE_GPMD,
						Code: gmp.GMPD_NODE_CONN_DENIED,
						Body: "Node not found",
					}
					respBts, err := json.Marshal(respMsg)
					if err != nil {
						fmt.Println("Error (connection.go:136): ", err)
					}
					conn.Write(respBts)
					continue
				}
				nodeBts, err := json.Marshal(respNode)
				if err != nil {
					fmt.Println("Error (connection.go:157): ", err)
				}
				respMsg := gmp.Message{
					Type: gmp.MESSAGE_TYPE_GPMD,
					Code: gmp.GMPD_NODE_CONN_ACCEPT,
					Body: string(nodeBts),
				}
				respBts, err := json.Marshal(respMsg)
				if err != nil {
					fmt.Println("Error (connection.go:166): ", err)
				}
				conn.Write(respBts)

			} else if msg.Code == gmp.GMPD_NODE_DISCOVER {
				items := []string{}
				for name, node := range nds {
					if node.Active {
						items = append(items, name)
					}

				}

				respMsg := gmp.Message{
					Type: gmp.MESSAGE_TYPE_GPMD,
					Code: gmp.GMPD_NODE_DISCOVER,
					Body: strings.Join(items, ","),
				}
				respBts, err := json.Marshal(respMsg)
				if err != nil {
					fmt.Println("Error (connection.go:182): ", err)
				}
				conn.Write(respBts)

			}
		} else if msg.Type == gmp.MESSAGE_TYPE_WASM_REQ {
			if msg.Code == gmp.WASM_REQ_MODULE {
				var req gmp.WasmRequest

				err := json.Unmarshal([]byte(msg.Body), &req)
				if err != nil {
					fmt.Printf("Error (connection.go:196): %v\n", err)
				}
				modRes, err := mod.GetModule(req.Module, req.Version)
				if err != nil {
					resMsg := gmp.Message{
						Type: gmp.MESSAGE_TYPE_WASM_REQ,
						Code: gmp.WASM_REQ_RESPONSE_NOT_FOUND,
					}

					modResBts, err := json.Marshal(resMsg)
					if err != nil {
						fmt.Printf("Error (connection.go:196): %v\n", err)
					}
					conn.Write(modResBts)
					continue
				}

				modResBts, err := json.Marshal(modRes)
				if err != nil {
					fmt.Printf("Error (connection.go:196): %v\n", err)
				}

				resMsg := gmp.Message{
					Type: gmp.MESSAGE_TYPE_WASM_REQ,
					Code: gmp.WASM_REQ_RESPONSE_W_BODY,
					Body: string(modResBts),
				}

				modResBts, err = json.Marshal(resMsg)
				if err != nil {
					fmt.Printf("Error (connection.go:196): %v\n", err)
				}

				_, err = conn.Write(modResBts)
				if err != nil {
					fmt.Println("Error: ", err)
				}

			}
		}

	}
}
