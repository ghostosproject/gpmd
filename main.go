package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"runtime"
	"slices"
	"syscall"

	"github.com/ghostosproject/gpmd/background"
	"github.com/ghostosproject/gpmd/module"
	"github.com/ghostosproject/gpmd/server"
)

func main() {
	uploadCmd := flag.NewFlagSet("upload", flag.ExitOnError)
	name := uploadCmd.String("name", "test", "Module Name")
	version := uploadCmd.String("version", "", "Module Version")
	file := uploadCmd.String("file", "", "Location of File")

	// Parse the flags
	flag.Parse()
	args := os.Args
	if len(args) <= 1 {
		fmt.Println("--help")
	}

	switch args[1] {
	case "start":
		if len(args) > 2 && args[2] == "nd" {
			mod := module.CreateModuleService()
			go background.Server(&mod)
			server.Server(&mod)
		} else {
			if isDetachedMode() {
				mod := module.CreateModuleService()
				// run a command server for running commands
				go background.Server(&mod)
				server.Server(&mod)
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
			uploadCmd.Parse(os.Args[3:])

			cn, err := net.Dial("tcp", "localhost:1111")
			if err != nil {
				fmt.Printf("Error connecting to the node\n")
			}

			bdy := fmt.Appendf(nil, "upload,%s,%s,%s", *name, *version, *file)

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
