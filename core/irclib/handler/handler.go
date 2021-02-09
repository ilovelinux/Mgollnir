package handler

import (
	"strings"

	"../../../core"
	"../commands"
)

func Handle(bot core.Bot, ircmsg core.IRCMsg) {
	switch ircmsg.Command {
	case "001": // RPL_WELCOME
		for _, channel := range bot.Server.Channels {
			bot.Server.Send(commands.Join(channel))
		}
	case "CTCP":
		handleCtcp(bot, ircmsg)
	case "JOIN":
		if bot.Identity.Username != ircmsg.User.Nickname {
			channel := strings.TrimPrefix(ircmsg.Parameters, ":")
			bot.Server.Send(commands.Privmsg(channel, ircmsg.User.Nickname+"!"))
		}
	case "PING":
		bot.Server.Send(commands.Pong(ircmsg.Parameters))
	case "PRIVMSG":
		handlePrivmsg(bot, ircmsg)
	}
}
