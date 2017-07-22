package commands

import (
	"strings"
	"unicode"

	"github.com/santiclause/eden/models"
)

var (
	commands      []*Command
	DefaultPrefix = "."
)

type Command struct {
	allowNoWhitespace bool
	command           string
	commandFunc       func(Message) string
	function          ExecuteFunc
	minArgs           int
	maxArgs           int
	permission        *models.Permission
	prefix            string
}

func (command *Command) Execute(message Message, context CommandContext) {
	prefix := command.prefix + command.command
	if command.commandFunc != nil {
		prefix = command.commandFunc(message)
	}
	if !strings.HasPrefix(message.Content, prefix) {
		return
	}
	remainder := strings.Replace(message.Content, prefix, "", 1)
	if len(remainder) > 0 && !command.allowNoWhitespace && remainder[0] != '\t' && remainder[0] != ' ' {
		return
	}
	args := parseArgs(remainder)
	if len(args) > command.maxArgs || len(args) < command.minArgs {
		return
	}
	if command.permission != nil && !context.Authorize(message.Source, *command.permission) {
		return
	}
	context.Execute(command.function, message, args...)
}

func ExecuteCommands(message Message, context CommandContext) {
	for _, command := range commands {
		command.Execute(message, context)
	}
}

func parseArgs(argstring string) []string {
	var args []string
	inQuotes := false
	escape := false
	arg := ""
	for _, c := range argstring {
		if inQuotes {
			if escape {
				if c == '"' {
					arg += string('"')
				} else if c == '\\' {
					arg += string('\\')
				} else {
					arg += string('\\') + string(c)
				}
				escape = false
			} else if c == '\\' {
				escape = true
			} else if c == '"' {
				inQuotes = false
			} else {
				arg += string(c)
			}
		} else {
			if escape {
				if unicode.IsSpace(c) {
					arg += string(c)
				} else if c == '"' {
					arg += string('"')
				} else {
					arg += string('\\') + string(c)
				}
				escape = false
			} else {
				if c == '\\' {
					escape = true
				} else if c == '"' {
					inQuotes = true
				} else if unicode.IsSpace(c) {
					if arg != "" {
						args = append(args, arg)
						arg = ""
					}
				} else {
					arg += string(c)
				}
			}
		}
	}
	if escape {
		arg += string('\\')
	}
	if arg != "" {
		args = append(args, arg)
	}
	return args
}

func NewCommand(command string, function ExecuteFunc, opts ...commandOption) (*Command, error) {
	c := &Command{
		command:  command,
		function: function,
		prefix:   DefaultPrefix,
	}
	for _, opt := range opts {
		err := opt(c)
		if err != nil {
			return nil, err
		}
	}
	commands = append(commands, c)
	return c, nil
}

type commandOption func(*Command) error

type ExecuteFunc func(CommandContext, Message, ...string)

type Message struct {
	Content string
	Source  User
	Public  bool
	Target  string
}

type User struct {
	Id          *int
	Name        string
	DisplayName string
}

func WithArgs(numArgs int) commandOption {
	return func(c *Command) error {
		c.minArgs = numArgs
		c.maxArgs = numArgs
		return nil
	}
}

func WithVarArgs(min, max int) commandOption {
	return func(c *Command) error {
		c.minArgs = min
		c.maxArgs = max
		return nil
	}
}

func WithAllowNoWhitespace() commandOption {
	return func(c *Command) error {
		c.allowNoWhitespace = true
		return nil
	}
}

func WithPermissionCheck(permission models.Permission) commandOption {
	return func(c *Command) error {
		c.permission = &permission
		return nil
	}
}

func WithPrefix(prefix string) commandOption {
	return func(c *Command) error {
		c.prefix = prefix
		return nil
	}
}

func WithCommandFunc(commandFunc func(Message) string) commandOption {
	return func(c *Command) error {
		c.commandFunc = commandFunc
		return nil
	}
}

type CommandContext interface {
	Execute(ExecuteFunc, Message, ...string)
	Authorize(User, models.Permission) bool
	SendToUser(User, string)
	SendToChannel(string, string)
}
