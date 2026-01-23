package command

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/cardvark/blog-aggregator/internal/database"
	"github.com/google/uuid"
)

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
