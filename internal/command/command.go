package command

import (
	"fmt"
	"os"

	"github.com/cardvark/blog-aggregator/internal/config"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		fmt.Println("Error: username required")
		os.Exit(1)
		return fmt.Errorf("Error: no arguments found for %s:\n", cmd.name)
	}

	userName := cmd.args[0]

	s.config.SetUser(userName)

	fmt.Printf("Username '%s' has been set.", userName)

	return nil
}

func (c commands) Run(s *state, cmd command) error {
	name := cmd.name
	cmdFunc, ok := c.commandMap[name]
	if !ok {
		return fmt.Errorf("No such command found: %s\n", name)
	}

	err := cmdFunc(s, cmd)
	if err != nil {
		return fmt.Errorf("Error running function '%s': %v", name, err)
	}

	return nil
}

func (c commands) register(name string, f func(*state, command) error) {
	c.commandMap[name] = f
}

func GetState(c config.Config) state {
	return state{config: &c}
}

func GetCommands() commands {
	cmds := commands{
		commandMap: make(map[string]func(*state, command) error),
	}

	cmds.register(
		"login",
		handlerLogin,
	)

	return cmds
}

func NewCommand(name string, args []string) command {
	return command{
		name: name,
		args: args,
	}
}
