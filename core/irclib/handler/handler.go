package handler

import (
	"../../../core"
	"../commands"
)

func Handle(sendq chan string, ircmsg core.IRCMsg) {
	switch ircmsg.Command {
	case "001":
		// TODO: send JOIN command
		break
	case "PING":
		sendq <- commands.Pong(ircmsg.Parameters)
	case "CTCP":
		handleCtcp(sendq, ircmsg)
	case "PRIVMSG":
		handlePrivmsg(sendq, ircmsg)
	}
}
