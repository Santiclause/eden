package main

import (
	irc "github.com/fluffle/goirc/client"
)

func blah() *irc.Conn {
	cfg := irc.NewConfig("nickname", "ident", "long name")
	cfg.Server = "server"
	cfg.Version = "versionnnn"
	cfg.QuitMessages = "so long gay bowser"
	return irc.Client(cfg)
}

func handlers(conn *irc.Conn) {
	conn.Handle("name", func(conn *irc.Conn, line *irc.Line) {
		switch line.Cmd {
		case "whatever":
			conn.Privmsg("someone", "something")
		default:
			conn.Privmsg("channel I guess", "something else")
		}
	})
}
