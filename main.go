package main

import (
	"fmt"
	"os"
)

func main() {
	prog := os.Args[0]
	if len(os.Args) != 3 {
		fmt.Printf("Usage: %s [listen address] [target]\n", prog)
		os.Exit(1)
	}
	args := os.Args[1:]

	listenAddress := args[0]
	targetAddress := args[1]

	RunTcpTunnel(listenAddress, targetAddress)
}