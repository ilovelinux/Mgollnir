package extra

import (
	"bufio"
	"os"
	"strings"

	"../irclib/commands"
)

func ConsoleReader(sendq chan string) {
	reader := bufio.NewReader(os.Stdin)
	channel := ""
	for {
		text, _ := reader.ReadString('\n')
		text = strings.TrimRight(text, "\n")

		switch text[0] {
		case '#':
			channel = text
		case '/':
			sendq <- text
		default:
			sendq <- commands.Privmsg(channel, text)
		}
	}
}
