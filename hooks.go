package main

// import irc "github.com/fluffle/goirc/client"

// type hookWrapper struct {
// 	event   string
// 	handler irc.Handler
// }

// var (
// 	hooks []*hookWrapper
// )

// // This only happens during init
// func addHook(name string, handler irc.Handler) {
// 	wrapper := &hookWrapper{
// 		event:   name,
// 		handler: handler,
// 	}
// 	hooks = append(hooks, wrapper)
// }

// func insertHooks(conn *irc.Conn) {
// 	for _, hook := range hooks {
// 		conn.Handle(hook.event, hook.handler)
// 	}
// }

// func test(c *irc.Conn, line *irc.Line) {
// 	line.Nick
// }

// type Authorizer interface {
// 	Authorized() bool
// }

import (
	"github.com/santiclause/eden/commands"
)

func init() {
	commands.NewCommand("hello", func(ctx commands.CommandContext, msg commands.Message, args ...string) {
		ctx.SendToChannel(msg.Target, "Hello world!")
	})
}
