package command

import (
	"github.com/cardvark/blog-aggregator/internal/config"
	"github.com/cardvark/blog-aggregator/internal/database"
)

type state struct {
	config *config.Config
	db     *database.Queries
}

type command struct {
	name string
	args []string
}

type commands struct {
	commandMap map[string]func(*state, command) error
}
