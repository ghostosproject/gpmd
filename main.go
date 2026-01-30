package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"slices"
	"strconv"
	"syscall"

	"github.com/ghostosproject/gpmd/module"
	"github.com/ghostosproject/gpmd/server"
)

// create the node structure
// what needs to be saved?

func main() {
	args := os.Args
	if len(args) <= 1 {
		fmt.Println("--help")
	}

	switch args[1] {
	case "start":
		if len(args) > 2 && args[2] == "nd" {
			module.CreateModuleService()
			server.Server()
		} else {
			if isDetachedMode() {
				module.CreateModuleService()
				server.Server()
				return
			}
			err := rerunDetached()
			if err != nil {
				log.Fatal(err)
			}

			fmt.Println("Main process exiting, detached process is running in the background.")
		}
	case "kill":
		bts, err := os.ReadFile("access.log")
		pid, err := strconv.Atoi(string(bts))
		process, err := os.FindProcess(pid)

		if err != nil {
			fmt.Printf("Error finding process: %v\n", err)
			return
		}
		process.Kill()
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
		// Unix-like systems (Linux, macOS): use Setpgid to create a new process group
		// and redirect I/O to /dev/null to survive terminal exit (SIGHUP)
		cmd.SysProcAttr = &syscall.SysProcAttr{
			Setpgid: true,
		}
		// Redirect standard files to /dev/null to prevent the parent terminal from holding the child
		cmd.Stdin = nil // or os.Open(os.DevNull)
		cmd.Stdout = nil
		cmd.Stderr = nil
	}

	err = cmd.Start()
	if err != nil {
		return err
	}

	pid := cmd.Process.Pid
	f, err := os.OpenFile("access.log", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	f.Write([]byte(strconv.Itoa(pid)))

	// Release the process handle in the parent process so it doesn't wait for the child
	err = cmd.Process.Release()
	if err != nil {
		log.Printf("Warning: Failed to release process: %v\n", err)
	}
	return nil
}
