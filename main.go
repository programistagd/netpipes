package main

import (
	"fmt"
	"os"
	"netpipes/netpipe"
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

	netpipe.RunTcpTunnel(listenAddress, targetAddress)
}