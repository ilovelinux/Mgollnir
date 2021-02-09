package handler

import (
	"strings"

	"../../../core"
	"../../driver"
	"../commands"
)

func Handle(d driver.Driver, ircmsg core.IRCMsg) {
	switch ircmsg.Command {
	case "001": // RPL_WELCOME
		for _, channel := range d.Server.Channels {
			d.Sendq <- commands.Join(channel)
		}
	case "CTCP":
		handleCtcp(d.Sendq, ircmsg)
	case "JOIN":
		if d.Identity.Username != ircmsg.User.Nickname {
			channel := strings.TrimPrefix(ircmsg.Parameters, ":")
			d.Sendq <- commands.Privmsg(channel, ircmsg.User.Nickname+"!")
		}
	case "PING":
		d.Sendq <- commands.Pong(ircmsg.Parameters)
	case "PRIVMSG":
		handlePrivmsg(d.Sendq, ircmsg)
	}
}
