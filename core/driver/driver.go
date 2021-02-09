package driver

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"../../core"

	"../irclib/commands"
)

type Driver struct {
	Server   core.Server
	Identity core.Identity

	conn  *net.Conn
	Sendq chan string
	Recvq chan string
	done  chan bool
}

func New(server core.Server, identity core.Identity) *Driver {
	var d = &Driver{}
	d.Server = server
	d.Identity = identity
	return d
}

func (d *Driver) Connect() {
	for d.conn == nil {
		var conn net.Conn
		var err error
		if d.Server.SSL {
			conn, err = tls.Dial("tcp", d.Server.String(), nil)
		} else {
			conn, err = net.Dial("tcp", d.Server.String())
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
		d.Sendq <- commands.Nick(d.Identity)
		d.Sendq <- commands.User(d.Identity)

		// Handle Ctrl+C
		c := make(chan os.Signal)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		go func() {
			<-c
			d.Disconnect()
			os.Exit(0)
		}()
	}
}

func (d *Driver) Disconnect(reason ...string) {
	if reason == nil {
		reason = []string{"Bug!"}
	}

	d.Sendq <- commands.Quit(reason[0])
	close(d.Sendq)
	<-d.done

	(*d.conn).Close()
	d.conn = nil
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
