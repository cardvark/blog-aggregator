package command

import (
	"context"
	"fmt"
	"os"

	"github.com/cardvark/blog-aggregator/internal/config"
	"github.com/cardvark/blog-aggregator/internal/database"
)

func middlewareLoggedIn(
	handler func(s *state, cmd command, user database.User) error,
) func(*state, command) error {

	return func(s *state, cmd command) error {
		userName := s.config.Current_user_name
		user, err := s.db.GetUser(
			context.Background(),
			userName,
		)
		if err != nil {
			fmt.Println("Current user not found in database.")
			os.Exit(1)
			return err
		}

		return handler(s, cmd, user)
	}
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
	cmds.register("agg", middlewareLoggedIn(handlerAgg))
	cmds.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	cmds.register("feeds", middlewareLoggedIn(handlerFeeds))
	cmds.register("follow", middlewareLoggedIn(handlerFollow))
	cmds.register("following", middlewareLoggedIn(handlerFollowing))
	cmds.register("unfollow", middlewareLoggedIn(handlerUnfollow))

	return cmds
}

func NewCommand(name string, args []string) command {
	return command{
		name: name,
		args: args,
	}
}
