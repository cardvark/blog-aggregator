package command

import (
	"fmt"

	"github.com/cardvark/blog-aggregator/internal/config"
	"github.com/cardvark/blog-aggregator/internal/database"
)

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

func GetState(c config.Config, dbq *database.Queries) state {
	return state{config: &c, db: dbq}
}

func GetCommands() commands {
	cmds := commands{
		commandMap: make(map[string]func(*state, command) error),
	}

	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerUsers)
	cmds.register("agg", handlerAgg)
	cmds.register("addfeed", handlerAddFeed)
	cmds.register("feeds", handlerFeeds)
	cmds.register("follow", handlerFollow)
	cmds.register("following", handlerFollowing)

	return cmds
}

func NewCommand(name string, args []string) command {
	return command{
		name: name,
		args: args,
	}
}
