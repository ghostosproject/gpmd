package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"runtime"
	"slices"
	"syscall"

	"github.com/caarlos0/env/v11"
	"github.com/ghostosproject/gpmd/config"
	"github.com/ghostosproject/gpmd/gpmd"
	"github.com/ghostosproject/gpmd/misc"
)

func main() {
	// uploadCmd := flag.NewFlagSet("upload", flag.ExitOnError)
	// name := uploadCmd.String("name", "test", "Module Name")
	// version := uploadCmd.String("version", "", "Module Version")
	// file := uploadCmd.String("file", "", "Location of File")

	// Parse the flags

	path := os.Getenv("GPMD_CONFIG_PATH")
	cfg := config.Config{}
	// check if config.json exists in local path
	if path != "" {
		// read file from path
		bts, err := os.ReadFile(path)
		if err != nil {
			panic(err)
		}
		err = json.Unmarshal(bts, &cfg)
	} else {
		// look for file at root location
		//
		if misc.FileExists("./config.json") {
			bts, err := os.ReadFile("./config.json")
			if err != nil {
				panic(err)
			}
			err = json.Unmarshal(bts, &cfg)
		} else {
			// env.Parse populates the struct fields from environment variables
			if err := env.Parse(&cfg); err != nil {
				log.Fatalf("failed to parse env: %+v", err)
			}
			cfg.Peers = []string{}
			bts, err := json.MarshalIndent(cfg, "", "	")
			if err != nil {
				panic(err)
			}
			err = os.WriteFile("./config.json", bts, 0777)
			if err != nil {
				panic(err)
			}
		}

	}
	flag.Parse()
	args := os.Args
	if len(args) <= 1 {
		fmt.Println("--help")
	}

	switch args[1] {
	case "start":
		if len(args) > 2 && args[2] == "nd" {
			node := gpmd.CreateGPMD(true, cfg.Ports.P2P, cfg.PrivateKeyPath)
			fmt.Printf("%s/p2p/%s\n", node.P2PHost.Addrs()[1], node.Id)
			node.Run(cfg)
		} else {
			if isDetachedMode() {
				node := gpmd.CreateGPMD(true, cfg.Ports.P2P, cfg.PrivateKeyPath)
				fmt.Printf("%s/p2p/%s\n", node.P2PHost.Addrs()[1], node.Id)
				node.Run(cfg)
				return
			}
			err := rerunDetached()
			if err != nil {
				log.Fatal(err)
			}

			fmt.Println("GPMD is running in the background.")
		}
	case "wasm":
		if args[2] == "upload" {
			// uploadCmd.Parse(os.Args[3:])

			// cn, err := net.Dial("tcp", "localhost:1111")
			// if err != nil {
			// 	fmt.Printf("Error connecting to the node\n")
			// }

			// bdy := fmt.Appendf(nil, "upload,%s,%s,%s", *name, *version, *file)

			// _, err = cn.Write(bdy)
			// if err != nil {
			// 	fmt.Println("Error: ", err)
			// }

			// buffer := make([]byte, 1024) // Create a buffer of 1024 bytes
			// n, err := cn.Read(buffer)
			// if err != nil {
			// 	fmt.Println(err)
			// }
			// fmt.Println(string(buffer[:n]))

		}
	case "kill":
		cn, err := net.Dial("tcp", "localhost:1111")
		if err != nil {
			fmt.Printf("Error connecting to the gpmd\n")
		}

		bdy := []byte("kill")

		_, err = cn.Write(bdy)
		if err != nil {
			fmt.Println("Error: ", err)
		}

		buffer := make([]byte, 1024) // Create a buffer of 1024 bytes
		n, err := cn.Read(buffer)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(string(buffer[:n]))
	}

}

func isDetachedMode() bool {
	return slices.Contains(os.Args, "--detached")
}

func rerunDetached() error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	// Prepare arguments for the new process, adding the "--detached" flag
	args := append(os.Args, "--detached")
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Dir = cwd

	// Configure SysProcAttr for proper detachment based on the OS
	if runtime.GOOS == "windows" {
		// Windows: use DETACHED_PROCESS flag
		// cmd.SysProcAttr = &syscall.SysProcAttr{
		// 	CreationFlags: syscall.CREATE_NEW_PROCESS_GROUP | syscall.DETACHED_PROCESS,
		// }
	} else {
		cmd.SysProcAttr = &syscall.SysProcAttr{
			Setpgid: true,
		}
		cmd.Stdin = nil // or os.Open(os.DevNull)
		cmd.Stdout = nil
		cmd.Stderr = nil
	}

	err = cmd.Start()
	if err != nil {
		return err
	}
	err = cmd.Process.Release()
	if err != nil {
		log.Printf("Warning: Failed to release process: %v\n", err)
	}
	return nil
}
