package main

import (
	"fmt"
	"os"

	"github.com/zic20/blog_aggregator/internal"
	"github.com/zic20/blog_aggregator/internal/config"
)

func main() {

	cfg, err := config.Read()
	if err != nil {
		fmt.Printf("Error loading config: %s", err)
		os.Exit(0)
	}

	state := internal.State{
		Config: &cfg,
	}

	commands := internal.Commands{
		Commands: make(map[string]func(*internal.State, internal.Command) error),
	}
	commands.Register("login", internal.HandlerLogin)

	args := os.Args

	if len(args) < 2 {
		fmt.Println("One command was expected")
		os.Exit(1)
	}
	command := internal.Command{
		Name: args[1],
		Args: args[2:],
	}

	commands.Run(&state, command)

	// cfg.SetUser("isaac")
	// updatedConfig, err := config.Read()
	// if err != nil {
	// 	fmt.Printf("Error loading config: %s", err)
	// 	os.Exit(0)
	// }

	// fmt.Print(updatedConfig)
}
