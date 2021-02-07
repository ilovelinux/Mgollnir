package handler

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"../../../core"
	"../../irclib"
	"../commands"
)

type ctcp struct{}

func handleCtcp(sendq chan string, ircmsg core.IRCMsg) {
	command, params := irclib.GetCommand(ircmsg.Parameters)

	response, err := runCTCPCommand(command, params)
	if err == nil {
		for _, line := range strings.Split(*response, "\n") {
			sendq <- commands.Notice(
				ircmsg.User.Nickname,
				fmt.Sprintf("\x01%s %s\x01", command, line),
			)
		}
	}
}

func runCTCPCommand(command, params string) (*string, error) {
	var c ctcp
	v := reflect.ValueOf(&c)
	method := v.MethodByName(strings.ToUpper(command))
	if method.IsValid() {
		in := []reflect.Value{reflect.ValueOf(params)}
		output := method.Call(in)[0].String()
		return &output, nil
	}
	return nil, fmt.Errorf("Method <%s %s> is not valid", command, params)
}

func (ctcp) CLIENTINFO(oarams string) string {
	var t ctcp
	val := reflect.ValueOf(&t)
	methods := []string{}
	for i := 0; i < val.NumMethod(); i++ {
		fnName := val.Type().Method(i).Name
		if fnName == "CLIENTINFO" {
			continue
		}
		methods = append(methods, fnName)
	}
	return strings.Join(methods, " ")
}

func (ctcp) FINGER(params string) string {
	return "FUCK YOU\n┌∩┐(◣_◢)┌∩┐"
}

func (ctcp) PING(params string) string {
	return params
}

func (ctcp) SOURCE(params string) string {
	return "https://github.com/ilovelinux/Mgollnir"
}

func (ctcp) TIME(params string) string {
	return time.Now().UTC().Format(time.RFC1123)
}

func (ctcp) USERINFO(params string) string {
	return "Mgollnir"
}

func (ctcp) VERSION(params string) string {
	return "8=========D"
}
