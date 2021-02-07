package main

import (
	"fmt"

	"./core"
	"./core/driver"
	"./core/extra"
	"./core/irclib"
	"./core/irclib/handler"
)

var server = core.Server{
	Server: "irc.azzurra.org",
	Port:   6667,
	SSL:    false,
	Channels: []core.Channel{
		core.Channel{Name: "#unity", Key: ""},
	},
}

var identity = core.Identity{
	Username:   "Mgollnir",
	Hostname:   "*",
	Servername: "*",
	Realname:   "Mgollnir",
}

func main() {
	d := driver.New(server, identity)
	d.Connect()
	go extra.ConsoleReader(d.Sendq)
	for {
		line := <-d.Recvq
		fmt.Println(line)
		ircmsg := irclib.Parser(line)
		handler.Handle(d.Sendq, ircmsg)
	}
}
