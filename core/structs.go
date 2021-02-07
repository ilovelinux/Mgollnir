package core

import "fmt"

type Server struct {
	Server string
	Port   int
	SSL    bool

	Channels []Channel
}

func (s Server) String() string {
	return fmt.Sprintf("%s:%d", s.Server, s.Port)
}

type Channel struct {
	Name string
	Key  string
}

type Identity struct {
	Username   string
	Hostname   string
	Servername string
	Realname   string
}

type User struct {
	Nickname string
	Identity string
	Hostname string
}

type IRCMsg struct {
	Command    string
	Parameters string

	User *User
}

type Privmsg struct {
	User    *User
	Channel string
	Text    string
}
