package main

import (
	"os"
	"netpipes/netpipe/tcp"
	"netpipes/netpipe/udp"
	"flag"
)

func main() {
	isUdpPtr := flag.Bool("udp", false, "Create UDP tunnel instead of TCP tunnel")
	listenAddrPtr := flag.String("from", "", "Address to listen on")
	targetAddrPtr := flag.String("to", "", "Address to redirect incoming connections to")

	flag.Parse()
	if *listenAddrPtr == "" || *targetAddrPtr == "" {
		flag.Usage()
		os.Exit(1)
	}

	if *isUdpPtr {
		udp.RunUdpTunnel(*listenAddrPtr, *targetAddrPtr)
	} else {
		tcp.RunTcpTunnel(*listenAddrPtr, *targetAddrPtr)
	}
}
