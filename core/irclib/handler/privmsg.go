package handler

import (
	"strings"

	"../../../core"
	"../../../plugins"
)

func handlePrivmsg(sendq chan string, ircmsg core.IRCMsg) {
	channel, text := func() (string, string) {
		x := strings.SplitN(ircmsg.Parameters, " ", 2)
		return x[0], x[1][1:]
	}()

	privmsg := core.Privmsg{
		User:    ircmsg.User,
		Channel: channel,
		Text:    text,
	}

	plugins.Dispatch(sendq, privmsg)
}
