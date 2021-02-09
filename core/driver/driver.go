package driver

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"time"
)

type Driver struct {
	addr string
	port int
	ssl  bool

	conn  *net.Conn
	Sendq chan string
	Recvq chan string
	done  chan bool
}

func New(addr string, port int, ssl bool) Driver {
	d := &Driver{}
	d.addr = addr
	d.port = port
	d.ssl = ssl

	return *d
}

func (d *Driver) Connect() {
	for d.conn == nil {
		var conn net.Conn
		var err error
		addr := fmt.Sprintf("%s:%d", d.addr, d.port)
		if d.ssl {
			conn, err = tls.Dial("tcp", addr, nil)
		} else {
			conn, err = net.Dial("tcp", addr)
		}
		if err != nil {
			log.Println(err)
			time.Sleep(5 * time.Second)
			continue
		}
		d.conn = &conn
		d.Sendq = make(chan string)
		d.Recvq = make(chan string)
		d.done = make(chan bool)
		go send(d)
		go recv(d)
	}
}

func (d *Driver) Disconnect() {
	if d.conn != nil {
		close(d.Sendq)
		<-d.done

		(*d.conn).Close()
		d.conn = nil
	}
}

func send(d *Driver) {
	for d.conn != nil {
		msg, ok := <-d.Sendq
		if ok {
			fmt.Println("Send:", msg)
			(*d.conn).Write(append([]byte(msg), '\n', '\r'))
		} else {
			d.done <- true
			return
		}
	}
}

func recv(d *Driver) {
	partialMessage := []byte{}
	for d.conn != nil {
		data := make([]byte, 4096)
		length, err := (*d.conn).Read(data)
		if err == nil {
			data = append(partialMessage, data[:length]...)
			partialMessage = []byte{}

			for _, message := range bytes.Split(data, []byte{'\n'}) {
				if bytes.HasSuffix(message, []byte{'\r'}) {
					d.Recvq <- string(message[:len(message)-1])
				} else {
					partialMessage = message
				}
			}
		}
	}
}
