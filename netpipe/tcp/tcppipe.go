package tcp

import (
	"net"
	"fmt"
	"os"
	"netpipes/netpipe/shared"
)

func RunTcpTunnel(listenAddress string, targetAddress string) {
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

type done struct {
	ch chan interface{}
}

func makeDone() *done {
	return &done{make(chan interface{})}
}

func fulfillDone(done *done) {
	// TODO there is potential for a race condition here
	if (done.ch != nil) {
		close(done.ch)
		done.ch = nil
	}
}

/*
Ties two connections so that all data send on one of them is redirected to the other.
Works with stream-like connections (TCP, not UDP).
 */
func streamTie(c1 net.Conn, c2 net.Conn) {
	done := makeDone()
	go redirect(c1, c2, done)
	redirect(c2, c1, done)
}

func redirect(from net.Conn, to net.Conn, done *done) {
	fc := shared.StartReader(from)
	tc := shared.StartWriter(to)
	defer close(tc)
	defer fulfillDone(done)

	for {
		select {
		case msg, ok := <-fc:
			if !ok {
				return
			} else {
				tc <- msg
			}
		case <-done.ch:
			return
		}
	}
}
