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

	listenAddress := args[0]
	targetAddress := args[1]

	fmt.Printf("Setting up a TCP tunnel redirecting all connections from %s to %s\n", listenAddress, targetAddress)

	l, err := net.Listen("tcp", listenAddress)
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

		go handleConnection(conn, targetAddress)
	}
}

func handleConnection(inbound net.Conn, targetAddress string) {
	defer inbound.Close()

	outbound, err := net.Dial("tcp", targetAddress)
	if err != nil {
		fmt.Println("Outbound connect error: ", err.Error())
		return
	}
	defer outbound.Close()

	fmt.Printf("Data incoming from %s is tunneled to %s\n", inbound.RemoteAddr().String(), outbound.RemoteAddr().String())

	streamTie(inbound, outbound)
}

// TODO currently releasing resources is not quite done (if one end ends the connection, the other may still be left open)

/*
Ties two connections so that all data send on one of them is redirected to the other.
Works with stream-like connections (TCP, not UDP).
 */
func streamTie(c1 net.Conn, c2 net.Conn) {
	go redirect(c1, c2)
	redirect(c2, c1)
}

type Message []byte

func reader(c net.Conn, ch chan Message) {
	defer close(ch)
	buff := make([]byte, 1024)
	for {
		length, err := c.Read(buff)
		if err != nil {
			fmt.Printf("Error reading: ", err.Error())
			return
		}
		tmp := make([]byte, length)
		copy(tmp, buff)
		ch <- tmp
	}
}

func startReader(c net.Conn) <-chan Message {
	ch := make(chan Message, 10)
	go reader(c, ch)
	return ch
}

func writer(ch chan Message, c net.Conn) {
	for msg := range ch {
		c.Write(msg)
	}
}

func startWriter(c net.Conn) chan<- Message {
	ch := make(chan Message, 10)
	go writer(ch, c)
	return ch
}

func redirect(from net.Conn, to net.Conn) {
	fc := startReader(from)
	tc := startWriter(to)

	for msg := range fc {
		tc <- msg
	}
}
