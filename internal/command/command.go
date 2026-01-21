package command

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/cardvark/blog-aggregator/internal/config"
	"github.com/cardvark/blog-aggregator/internal/database"
	"github.com/cardvark/blog-aggregator/internal/rss"
	"github.com/google/uuid"
)

func handlerAddFeed(s *state, cmd command) error {
	if len(cmd.args) < 2 {
		fmt.Println("Error: feed (1) name and (2) URL required.")
		os.Exit(1)
		return fmt.Errorf("Error: no arguments found for %s:\n", cmd.name)
	}

	feedName := cmd.args[0]
	feedURL := cmd.args[1]

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

	user_id := user.ID

	now := time.Now()
	var nullableTime sql.NullTime
	nullableTime.Time = now
	nullableTime.Valid = true

	feedParams := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: nullableTime,
		Name:      feedName,
		Url:       feedURL,
		UserID:    user_id,
	}

	feed, err := s.db.CreateFeed(
		context.Background(),
		feedParams,
	)
	if err != nil {
		fmt.Printf("Error creating feed entry: %v\n", err)
		return err
	}

	fmt.Printf("%#", feed)
	return nil
}

func handlerAgg(s *state, cmd command) error {

	feedURL := "https://www.wagslane.dev/index.xml"

	rssFeed, err := rss.FetchFeed(
		context.Background(),
		feedURL,
	)
	if err != nil {
		return err
	}

	fmt.Print(rssFeed)

	return nil
}

func handlerUsers(s *state, cmd command) error {
	users, err := s.db.GetUsers(
		context.Background(),
	)
	if err != nil {
		fmt.Printf("Error retrieving users: %v", err)
		os.Exit(1)
		return err
	}

	for _, user := range users {
		if s.config.Current_user_name == user.Name {
			fmt.Printf("* %s (current)\n", user.Name)
		} else {
			fmt.Printf("* %s\n", user.Name)
		}
	}

	return nil
}

func handlerReset(s *state, cmd command) error {
	err := s.db.DeleteUsers(
		context.Background(),
	)

	if err != nil {
		fmt.Printf("Error deleting users: %v\n", err)
		os.Exit(1)
		return err
	}

	fmt.Println("Deleted all users from 'users' table.")

	return nil

}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		fmt.Println("Error: username required")
		os.Exit(1)
		return fmt.Errorf("Error: no arguments found for %s:\n", cmd.name)
	}

	userName := cmd.args[0]

	// check if user is in DB:
	_, err := s.db.GetUser(
		context.Background(),
		userName,
	)

	if err != nil {
		fmt.Println("No such user found. Unable to log in.")
		os.Exit(1)
		return err
	}

	s.config.SetUser(userName)

	fmt.Printf("Username '%s' has been set.", userName)

	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		fmt.Println("Error: name required")
		return fmt.Errorf("Error: no arguments found for %s:\n", cmd.name)
	}

	name := cmd.args[0]
	now := time.Now()
	var nullableTime sql.NullTime
	nullableTime.Time = now
	nullableTime.Valid = true

	data := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: now,
		UpdatedAt: nullableTime,
		Name:      name,
	}

	user, err := s.db.CreateUser(
		context.Background(),
		data,
	)

	if err != nil {
		fmt.Printf("Error creating user: %v", err)
		os.Exit(1)
	}
	s.config.SetUser(name)

	fmt.Printf("User (%s) created. user data: %v", name, user)

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

	return cmds
}

func NewCommand(name string, args []string) command {
	return command{
		name: name,
		args: args,
	}
}
