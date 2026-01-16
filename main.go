package main

import (
	"fmt"
	"os"

	config "github.com/cardvark/blog-aggregator/internal/config"
)

func main() {
	homePath, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}

	test, err := config.Read(homePath)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%#v", test)
}
