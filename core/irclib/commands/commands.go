package commands

import (
	"fmt"
	"strings"

	core "../.."
)

func User(identity core.Identity) string {
	return fmt.Sprintf("USER %s %s %s :%s", identity.Username, identity.Hostname, identity.Servername, identity.Realname)
}

func Nick(identity core.Identity) string {
	return fmt.Sprintf("NICK %s", identity.Username)
}

// func MultipleJoin(channels []core.Channel) string {
// 	return fmt.Sprintf("JOIN %s %s", strings.Join(channels, ','), channel.key)
// }

func Join(channel core.Channel) string {
	return fmt.Sprintf("JOIN %s %s", channel.Name, channel.Key)
}

func Privmsg(receiver, message string) string {
	return fmt.Sprintf("PRIVMSG %s :%s", receiver, message)
}

func Notice(receiver, message string) string {
	return fmt.Sprintf("NOTICE %s :%s", receiver, message)
}

func Pong(daemons ...string) string {
	return fmt.Sprintf("PONG %s", strings.Join(daemons, " "))
}

func Quit(reason string) string {
	return fmt.Sprintf("QUIT %s", reason)
}
