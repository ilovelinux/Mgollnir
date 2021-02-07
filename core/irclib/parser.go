package irclib

import (
	"regexp"
	"strings"

	"../../core"
)

var senderRe = regexp.MustCompile(`^(\S+)!(\S+)@(\S+)$`)
var ctcpRe = regexp.MustCompile(`\x01([^\x01]+)\x01$`)

func splitHostmask(hostmask string) *core.User {
	m := senderRe.FindStringSubmatch(hostmask)
	if len(m) == 4 {
		return &core.User{
			Nickname: m[1],
			Identity: m[2],
			Hostname: m[3],
		}
	}
	return nil
}

func split2(s, sep string) (string, string) {
	x := strings.SplitN(s, sep, 2)
	return x[0], x[1]
}

// Parse a raw IRC message
func Parser(line string) core.IRCMsg {
	ircmsg := core.IRCMsg{}

	first, rest := split2(line, " ")
	if strings.HasPrefix(first, ":") {
		ircmsg.User = splitHostmask(first[1:])
		ircmsg.Command, ircmsg.Parameters = split2(rest, " ")
		if ircmsg.Command == "PRIVMSG" {
			// Handle CTCP messages
			if ctcpRe.MatchString(ircmsg.Parameters) {
				ircmsg.Command = "CTCP"
				ircmsg.Parameters = ctcpRe.FindStringSubmatch(ircmsg.Parameters)[1]
			}
		}
	} else {
		ircmsg.Command = first
		ircmsg.Parameters = rest
	}

	ircmsg.Command = strings.ToUpper(ircmsg.Command)

	return ircmsg
}
