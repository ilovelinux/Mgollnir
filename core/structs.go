package core

import (
	"./driver"
)

type Bot struct {
	Server   Server
	Identity Identity
}

type Server struct {
	Server string
	Port   int
	SSL    bool

	Driver driver.Driver

	Channels []Channel
}

func (s *Server) Connect() {
	s.Driver = driver.New(s.Server, s.Port, s.SSL)
	s.Driver.Connect()
}

func (s *Server) Send(l string) {
	s.Driver.Sendq <- l
}

func (s *Server) Recv() string {
	return <-s.Driver.Recvq
}

func (s *Server) Disconnect() {
	s.Driver.Disconnect()
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
