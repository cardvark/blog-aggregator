package main

import (
	"fmt"
	"os"

	"github.com/cardvark/blog-aggregator/internal/command"
	"github.com/cardvark/blog-aggregator/internal/config"
)

func main() {
	args := os.Args

	if len(args) < 2 {
		fmt.Println("Please include a command.")
		os.Exit(1)
	}

	homePath, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}

	config.InitPaths(homePath)

	cfg, err := config.Read()
	if err != nil {
		fmt.Println(err)
	}

	cfgState := command.GetState(cfg)

	cmdName := args[1]
	var cmdArgs []string

	if len(args) > 2 {
		cmdArgs = args[2:]
	}

	newCommand := command.NewCommand(cmdName, cmdArgs)
	commands := command.GetCommands()

	err = commands.Run(&cfgState, newCommand)
	if err != nil {
		fmt.Println(err)
	}

}
