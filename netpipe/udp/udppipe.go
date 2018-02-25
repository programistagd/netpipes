package udp

import (
	"fmt"
	"net"
	"os"
	"sync"
	"netpipes/netpipe/shared"
	"time"
)

var connections sync.Map

func RunUdpTunnel(listenAddress string, targetAddress string) {
	fmt.Printf("Setting up a UDP tunnel redirecting all connections from %s to %s\n", listenAddress, targetAddress)
	l, err := net.ListenPacket("udp", listenAddress)
	if err != nil {
		fmt.Println("Listen error: ", err.Error())
		os.Exit(1)
	}
	defer l.Close()

	buff := make([]byte, 1024)
	for {
		n, addr, err := l.ReadFrom(buff)
		if err != nil {
			fmt.Println("Receive error: ", err.Error())
			continue
		}

		tmp := make([]byte, n)
		copy(tmp, buff)

		c, found := connections.Load(addr.String())

		var conn net.Conn
		if found {
			conn = c.(net.Conn)
		} else {
			conn, err = net.Dial("udp", targetAddress)
			if err != nil {
				fmt.Println("Dial error: ", err.Error())
				continue
			}
			connections.Store(addr.String(), conn)
			go redirect(conn, addr, l)
			fmt.Printf("Data incoming from %s is tunneled to %s\n", addr.String(), conn.RemoteAddr().String())
		}
		conn.Write(tmp)
	}

}

func redirect(from net.Conn, fromAddr net.Addr, to net.PacketConn) {
	defer from.Close()
	defer connections.Delete(fromAddr.String())

	inc := shared.StartReader(from)
	for {
		select {
		case msg, ok := <-inc:
			if !ok {
				return
			} else {
				to.WriteTo(msg, fromAddr)
			}
		case <-time.After(300 * time.Second):
			fmt.Printf("%s timed out\n", fromAddr.String())
			return
		}
	}
}