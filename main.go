package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	prog := os.Args[0]
	if len(os.Args) != 3 {
		fmt.Printf("Usage: %s [listen address] [target]\n", prog)
		os.Exit(1)
	}
	args := os.Args[1:]

	listen_address := args[0]
	target_address := args[1]

	fmt.Printf("Setting up a TCP tunnel redirecting all connections from %s to %s\n", listen_address, target_address)

	l, err := net.Listen("tcp", listen_address)
	if err != nil {
		fmt.Println("Listen error: ", err.Error())
		os.Exit(1)
	}
	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Accept error: ", err.Error())
			continue
		}

		go handleConnection(conn, target_address)
	}
}

func handleConnection(inbound net.Conn, target_address string) {
	defer inbound.Close()

	outbound, err := net.Dial("tcp", target_address)
	if err != nil {
		fmt.Println("Outbound connect error: ", err.Error())
		return
	}
	defer outbound.Close()

	fmt.Printf("Data incoming from %s is tunneled to %s\n", inbound.RemoteAddr().String(), outbound.RemoteAddr().String())

	/*buf := make([]byte, 1024)
	for {
		len, err := inbound.Read(buf)
		if err != nil {

		}
	}*/
}