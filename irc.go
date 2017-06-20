package main

import (
	"log"
	"sync"

	irc "github.com/fluffle/goirc/client"
	"github.com/santiclause/eden/commands"
	"github.com/santiclause/eden/models"
)

type IrcConn struct {
	// This is a map of nicknames to Eden Users. We store this to cache
	// nickserv lookups.
	users userMap
	conn  *irc.Conn
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
		users: makeMap(),
		conn:  irc.Client(cfg),
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
		c.users.remove(line.Nick)
	}
}

func (c *IrcConn) part() irc.HandlerFunc {
	return func(conn *irc.Conn, line *irc.Line) {
		c.users.remove(line.Nick)
	}
}

func (c *IrcConn) nick() irc.HandlerFunc {
	return func(conn *irc.Conn, line *irc.Line) {
		c.users.remove(line.Nick)
	}
}

type userMap struct {
	mapping map[string]*models.User
	sync.RWMutex
}

func makeMap() (m userMap) {
	m.mapping = make(map[string]*models.User)
	return
}

func (m *userMap) get(key string) (user *models.User, ok bool) {
	m.RLock()
	defer m.RUnlock()
	user, ok = m.mapping[key]
	return
}

func (m *userMap) set(key string, value *models.User) {
	m.Lock()
	defer m.Unlock()
	m.mapping[key] = value
}

func (m *userMap) remove(key string) {
	m.Lock()
	defer m.Unlock()
	delete(m.mapping, key)
}

// CommandContext interface methods

func (c *IrcConn) Execute(f commands.ExecuteFunc, message commands.Message, args []string) {
	f(message, args)
}

func (c *IrcConn) Authorize(userInfo commands.User, permission models.Permission) bool {
	user, ok := c.users.get(userInfo.Name)
	if !ok {
		//TODO query nickserv here
		// c.conn.Privmsgf("NickServ", "STATUS %s", userInfo.Name)
		// channel with timeout, etc
		return false
	}
	if user == nil {
		ircUser := models.IrcUser{
			Nickname: userInfo.Name,
		}
		if db.Where(&ircUser).First(&ircUser).RecordNotFound() {
			return false
		}
		user = new(models.User)
		if err := db.Model(&ircUser).Related(user).Error; err != nil {
			if config.DebugLevel("warning") {
				log.Printf("Error fetching user for ircUser: %s\n", err)
			}
			return false
		}
		c.users.set(userInfo.Name, user)
	}
	if err := user.GetPermissions(db); err != nil {
		if config.DebugLevel("warning") {
			log.Printf("Error fetching user permissions: %s\n", err)
		}
		return false
	}
	for _, p := range user.Permissions {
		if p == permission {
			return true
		}
	}
	return false
}

type ircOption func(*irc.Config)

func WithNickname(nickname string) ircOption {
	return func(cfg *irc.Config) {
		cfg.Me.Nick = nickname
	}
}
