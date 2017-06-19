package main

import (
	irc "github.com/fluffle/goirc/client"
	"github.com/santiclause/eden/commands"
)

type IrcConn struct {
	// This is a map of nicknames to Eden User IDs. We store this to cache
	// nickserv lookups.
	userMap map[string]int
	conn    *irc.Conn
}

var handlers = map[string]func(*IrcConn) irc.HandlerFunc{
	irc.QUIT:    (*IrcConn).quit,
	irc.PART:    (*IrcConn).part,
	irc.NICK:    (*IrcConn).nick,
	irc.PRIVMSG: (*IrcConn).commandHook,
}

func Connect(server string, opts ...ircOption) *IrcConn {
	cfg := irc.NewConfig(config.IrcNickname, config.IrcIdent, config.IrcName)
	cfg.Server = server
	cfg.Version = config.Version
	cfg.QuitMessage = config.IrcQuitMessage
	for _, opt := range opts {
		opt(cfg)
	}
	conn := &IrcConn{
		userMap: make(map[string]int),
		conn:    irc.Client(cfg),
	}
	for event, hook := range handlers {
		conn.conn.HandleFunc(event, hook(conn))
	}
	return conn
}

func (c *IrcConn) commandHook() irc.HandlerFunc {
	return func(conn *irc.Conn, line *irc.Line) {
		message := commands.Message{
			Content: line.Text(),
			Public:  line.Public(),
		}
		commands.ExecuteCommands(message, c)
	}
}

func (c *IrcConn) quit() irc.HandlerFunc {
	return func(conn *irc.Conn, line *irc.Line) {
		delete(c.userMap, line.Nick)
	}
}

func (c *IrcConn) part() irc.HandlerFunc {
	return func(conn *irc.Conn, line *irc.Line) {
		delete(c.userMap, line.Nick)
	}
}

func (c *IrcConn) nick() irc.HandlerFunc {
	return func(conn *irc.Conn, line *irc.Line) {
		delete(c.userMap, line.Nick)
	}
}

// CommandContext interface methods

func (c *IrcConn) Execute(f commands.ExecuteFunc, message commands.Message, args []string) {
	f(message, args)
}

func (c *IrcConn) Authorize(user commands.User, permission commands.Permission) bool {
	id, ok := c.userMap[user.Name]
	if !ok {
		//TODO get user id here
		id = id
	}
	// get sql user here
	return false
}

type ircOption func(*irc.Config)

func WithNickname(nickname string) ircOption {
	return func(cfg *irc.Config) {
		cfg.Me.Nick = nickname
	}
}
