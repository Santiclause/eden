package main

import irc "github.com/fluffle/goirc/client"

func init() {
	// addHook(irc.PRIVMSG, favesHook)
}

func favesHook(conn *irc.Conn, line *irc.Line) {
}
