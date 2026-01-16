package main

import (
	"fmt"
	"os"

	"github.com/cardvark/blog-aggregator/internal/config"
)

func main() {
	homePath, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}

	config.InitPaths(homePath)

	cfg, err := config.Read()
	if err != nil {
		fmt.Println(err)
	}

	cfg.SetUser("james")

	cfg, err = config.Read()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("%#v", cfg)
}
