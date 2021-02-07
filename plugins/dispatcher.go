package plugins

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"../core"
	"../core/irclib"
	"../core/irclib/commands"
)

const commandPrefix = "}"

var nestedCommandRe = regexp.MustCompile("(?:^|`)([^ `]+)( ?[^`]*)`")
var validCmdChars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"

type botCommand struct{}

func validateCommand(cmd string) string {
	if cmd == "" || strings.ContainsAny(cmd[:1], "0123456789_") {
		return ""
	}
	return strings.TrimFunc(cmd, func(r rune) bool {
		return !strings.ContainsRune(validCmdChars, r)
	})
}

func Dispatch(sendq chan string, privmsg core.Privmsg) {
	message := privmsg.Text
	result := ""
	cycles := 0

	// Handle nested commands.
	// User: `echo test`
	//  Bot: test
	for loc := []int{}; ; message = message[loc[1]:] {
		loc = nestedCommandRe.FindStringIndex(message)
		if loc == nil {
			result += message
			break
		}
		result += message[:loc[0]]
		command, args := func() (string, string) {
			x := nestedCommandRe.FindStringSubmatch(message[loc[0]:loc[1]])
			return strings.Title(x[1]), x[2]
		}()

		// Needed to avoid empty arguments. (E.g. `echo `)
		// It still accepts commands without arguments (E.g. `ping`)
		if args != " " {
			args = strings.TrimPrefix(args, " ")

			output, err := runBotCommand(command, args)
			if err == nil {
				result += *output
				cycles++
			} else {
				result += message[loc[0] : loc[1]-1]
			}
		} else {
			result += message[loc[0] : loc[1]-1]
		}
	}
	message = result

	// Handle commands by prefix.
	// User: }echo test
	//  Bot: test
	if strings.HasPrefix(message, commandPrefix) {
		command, args := irclib.GetCommand(message[1:])
		output, err := runBotCommand(command, args)
		if err == nil {
			result = *output
			cycles++
		} else {
			result = fmt.Sprintf("Invalid command <%s %s>", command, args)
		}
	}

	if cycles > 0 {
		result = fmt.Sprintf("%s (%d cycles)", result, cycles)
		for _, line := range strings.Split(result, "\n") {
			sendq <- commands.Privmsg(privmsg.Channel, line)
		}
	}
}

func runBotCommand(command, args string) (*string, error) {
	var botcmd botCommand
	v := reflect.ValueOf(&botcmd)
	method := v.MethodByName(strings.Title(command))
	if method.IsValid() {
		in := []reflect.Value{reflect.ValueOf(args)}
		output := method.Call(in)[0].String()
		return &output, nil
	}
	return nil, fmt.Errorf("Method <%s %s> is not valid", command, args)
}
