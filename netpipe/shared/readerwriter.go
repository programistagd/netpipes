package shared

import (
	"net"
	"fmt"
)

type Message []byte

func reader(c net.Conn, ch chan Message) {
	defer close(ch)
	buff := make([]byte, 1024)
	for {
		length, err := c.Read(buff)
		if err != nil {
			fmt.Println("Error reading: ", err.Error())
			return
		}
		tmp := make([]byte, length)
		copy(tmp, buff)
		ch <- tmp
	}
}

func StartReader(c net.Conn) <-chan Message {
	ch := make(chan Message, 10)
	go reader(c, ch)
	return ch
}

func writer(ch chan Message, c net.Conn) {
	for msg := range ch {
		c.Write(msg)
	}
}

func StartWriter(c net.Conn) chan<- Message {
	ch := make(chan Message, 10)
	go writer(ch, c)
	return ch
}