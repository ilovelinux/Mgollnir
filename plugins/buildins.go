package plugins

import (
	"fmt"
	"math/rand"
	"strings"
)

func (botCommand) Echo(args string) string {
	return args
}

func (botCommand) Ping(args string) string {
	return "Pong!"
}

func (botCommand) Ratto(args string) string {
	var mult int
	switch {
	case strings.Contains(args, "giga"), strings.Contains(args, "iper"):
		mult = 4
		if rand.Int()%10 == 0 {
			mult <<= 1
		}
	case strings.Contains(args, "mega"), strings.Contains(args, "maxi"):
		mult = 3
	case strings.Contains(args, "super"):
		mult = 2
	default:
		mult = 1
	}
	rat := "8=%s==D"
	if rand.Int()%10 >= 7 {
		rat += "O:"
	}
	mast := []string{"==", "===", "====", "======", "=================="}[rand.Int()%5]
	return fmt.Sprintf(rat, strings.Repeat(mast, mult))
}

func (botCommand) Reverse(args string) string {
	runes := []rune(args)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func (botCommand) Source(args string) string {
	return "https://github.com/ilovelinux/Mgollnir"
}
