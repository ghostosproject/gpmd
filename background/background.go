package background

import (
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/ghostosproject/gpmd/module"
)

func HandleBackgroundConnection(conn net.Conn, mod *module.ModuleService) {
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
		args := strings.Split(string(buffer[:n]), ",")
		if len(args) == 1 {
			if args[0] == "kill" {
				conn.Write([]byte("gmpd stopped"))
				conn.Close()
				os.Exit(0)
				return
			}
		}
		if len(args) != 4 {
			conn.Write([]byte("Not enough args"))
			conn.Close()
			return
		}
		if args[0] == "upload" {
			err = mod.SetModule(args[1], args[2], args[3])
			if err != nil {
				conn.Write([]byte(err.Error()))
				return
			}
			conn.Write([]byte("Module Uploaded Successfully!"))
			conn.Close()
			return
		}

	}
}
