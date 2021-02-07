package irclib

import (
	"strings"
)

func GetCommand(s string) (string, string) {
	x := append(strings.SplitN(s, " ", 2), "")
	return x[0], x[1]
}
