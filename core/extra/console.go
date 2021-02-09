package extra

import (
	"bufio"
	"os"
	"strings"

	"../../core"
	"../irclib/commands"
)

func ConsoleReader(server core.Server) {
	reader := bufio.NewReader(os.Stdin)
	channel := ""
	for {
		text, _ := reader.ReadString('\n')
		text = strings.TrimRight(text, "\n")

		switch text[0] {
		case '#':
			channel = text
		case '/':
			server.Send(strings.TrimPrefix(text, "/"))
		default:
			server.Send(commands.Privmsg(channel, text))
		}
	}
}
