package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"./core"
	"./core/extra"
	"./core/irclib"
	"./core/irclib/commands"
	"./core/irclib/handler"
)

var bot = core.Bot{
	Server: core.Server{
		Server: "irc.azzurra.org",
		Port:   6667,
		SSL:    false,
		Channels: []core.Channel{
			core.Channel{Name: "#unity", Key: ""},
		},
	},

	Identity: core.Identity{
		Username:   "Mgollnir",
		Hostname:   "*",
		Servername: "*",
		Realname:   "Mgollnir",
	},
}

func main() {
	bot.Server.Connect()

	bot.Server.Send(commands.Nick(bot.Identity))
	bot.Server.Send(commands.User(bot.Identity))

	// Handle Ctrl+C
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		bot.Server.Send(commands.Quit("Bug!"))
		os.Exit(0)
	}()

	go extra.ConsoleReader(bot.Server)
	for {
		line := bot.Server.Recv()
		fmt.Println(line)
		ircmsg := irclib.Parser(line)
		handler.Handle(bot, ircmsg)
	}
}
