package main

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"sync"
	"time"

	irc "github.com/fluffle/goirc/client"
	"github.com/santiclause/eden/commands"
	"github.com/santiclause/eden/models"
)

type IrcConn struct {
	// This is a map of nicknames to Eden Users. We store this to cache
	// nickserv lookups.
	users            userMap
	conn             *irc.Conn
	cfg              *irc.Config
	autojoinChannels []string
	nickservPassword string
	nickservTimeout  time.Duration
}

var handlers = map[string]func(*IrcConn) irc.HandlerFunc{
	irc.CONNECTED: (*IrcConn).connected,
	irc.MODE:      (*IrcConn).mode,
	irc.NICK:      (*IrcConn).nick,
	irc.PART:      (*IrcConn).part,
	irc.PRIVMSG:   (*IrcConn).commandHook,
	irc.QUIT:      (*IrcConn).quit,
}

// Connects to an IRC server with the given options.
func Connect(server, nickname string, opts ...ircOption) (*IrcConn, error) {
	cfg := irc.NewConfig(nickname)
	cfg.Server = server
	conn := &IrcConn{
		cfg:   cfg,
		users: makeMap(),
	}
	for _, opt := range opts {
		opt(conn)
	}
	conn.conn = irc.Client(cfg)
	for event, hook := range handlers {
		conn.conn.HandleFunc(event, hook(conn))
	}
	if err := conn.conn.Connect(); err != nil {
		return nil, err
	}
	return conn, nil
}

func (c *IrcConn) commandHook() irc.HandlerFunc {
	return func(conn *irc.Conn, line *irc.Line) {
		message := commands.Message{
			Content: line.Text(),
			Public:  line.Public(),
			Source: commands.User{
				Name: line.Nick,
			},
			Target: line.Target(),
		}
		commands.ExecuteCommands(message, c)
	}
}

func (c *IrcConn) connected() irc.HandlerFunc {
	return func(conn *irc.Conn, line *irc.Line) {
		if c.nickservPassword != "" {
			c.conn.Privmsgf("NickServ", "IDENTIFY %s", c.nickservPassword)
		} else {
			c.Autojoin()
		}
	}
}

func (c *IrcConn) mode() irc.HandlerFunc {
	return func(conn *irc.Conn, line *irc.Line) {
		if line.Args[0] == conn.Me().Nick && line.Args[1] == "+r" {
			c.Autojoin()
		}
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

func (c *IrcConn) Execute(f commands.ExecuteFunc, message commands.Message, args ...string) {
	f(c, message, args...)
}

func (c *IrcConn) Authorize(userInfo commands.User, permission models.Permission) bool {
	user, ok := c.users.get(userInfo.Name)

	// Cache miss, we don't have any information about this user
	if !ok {
		timeout := time.After(c.nickservTimeout)
		wait := make(chan bool, 1)
		remover := c.conn.HandleFunc(irc.PRIVMSG, func(conn *irc.Conn, line *irc.Line) {
			if line.Nick == "NickServ" && !line.Public() {
				if ok, _ := regexp.MatchString(fmt.Sprintf("^STATUS %s \\d$", regexp.QuoteMeta(userInfo.Name)), line.Text()); ok {
					if strings.HasSuffix(line.Text(), "3") {
						wait <- true
					} else {
						wait <- false
					}
				}
			}
		})
		c.conn.Privmsgf("NickServ", "STATUS %s", userInfo.Name)
		select {
		case ok = <-wait:
		case <-timeout:
		}
		remover.Remove()
		if ok {
			// This user is registered and identified with NickServ, so we want to
			// at least cache that this user is verified by NickServ.
			c.users.set(userInfo.Name, nil)
		} else {
			return false
		}
	}

	// Cache hit (or successful check with NickServ), but we don't have an Eden user for them yet.
	if user == nil {
		ircUser := models.IrcUser{
			Nickname: userInfo.Name,
		}
		if db.Where(&ircUser).First(&ircUser).RecordNotFound() {
			return false
		}
		user = new(models.User)
		if err := db.Model(&ircUser).Related(user).Error; err != nil {
			log.Printf("Error fetching user for ircUser: %s\n", err)
			return false
		}
		// Cache the Eden user
		c.users.set(userInfo.Name, user)
	}

	if err := user.GetPermissions(db); err != nil {
		log.Printf("Error fetching user permissions: %s\n", err)
		return false
	}
	for _, p := range user.Permissions {
		if p == permission {
			return true
		}
	}
	return false
}

func (c *IrcConn) SendToUser(userInfo commands.User, message string) {
	c.conn.Privmsg(userInfo.Name, message)
}

func (c *IrcConn) SendToChannel(channel, message string) {
	c.conn.Privmsg(channel, message)
}

// end interface definitions

func (c *IrcConn) Autojoin() {
	if c.autojoinChannels == nil {
		return
	}
	for _, channel := range c.autojoinChannels {
		c.conn.Join(channel)
	}
}

type ircOption func(*IrcConn)

func WithAutojoinChannels(channels []string) ircOption {
	return func(c *IrcConn) {
		c.autojoinChannels = channels
	}
}

func WithIdent(ident string) ircOption {
	return func(c *IrcConn) {
		if ident != "" {
			c.cfg.Me.Ident = ident
		}
	}
}

func WithName(name string) ircOption {
	return func(c *IrcConn) {
		if name != "" {
			c.cfg.Me.Name = name
		}
	}
}

func WithNickservPassword(password string) ircOption {
	return func(c *IrcConn) {
		c.nickservPassword = password
	}
}

func WithNickservTimeout(timeout time.Duration) ircOption {
	return func(c *IrcConn) {
		c.nickservTimeout = timeout
	}
}

func WithQuitMessage(quitMessage string) ircOption {
	return func(c *IrcConn) {
		if quitMessage != "" {
			c.cfg.QuitMessage = quitMessage
		}
	}
}

func WithTimeout(timeout time.Duration) ircOption {
	return func(c *IrcConn) {
		c.cfg.Timeout = timeout
	}
}

func WithVersion(version string) ircOption {
	return func(c *IrcConn) {
		if version != "" {
			c.cfg.Version = version
		}
	}
}
