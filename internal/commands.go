package internal

import (
	"errors"
	"fmt"
	"os"

	"github.com/zic20/blog_aggregator/internal/config"
)

type State struct {
	Config *config.Config
}
type Command struct {
	Name string
	Args []string
}
type Commands struct {
	Commands map[string]func(*State, Command) error
}

func (c *Commands) Run(s *State, cmd Command) error {
	handler, ok := c.Commands[cmd.Name]
	if !ok {
		return fmt.Errorf("%s command does not exist", cmd.Name)
	}

	err := handler(s, cmd)
	if err != nil {
		return fmt.Errorf("Error occured while running %s command", cmd.Name)
	}

	return nil
}

func (c *Commands) Register(name string, f func(*State, Command) error) {
	c.Commands[name] = f
}

func HandlerLogin(s *State, cmd Command) error {
	if len(cmd.Args) == 0 {
		fmt.Println("username is required")
		os.Exit(1)
		return errors.New("login handler expects one argument")
	}

	if err := s.Config.SetUser(cmd.Args[0]); err != nil {
		return err
	}
	fmt.Printf("User %s has been set", cmd.Args[0])
	return nil
}

// func GetCommands() map[string]cliCommand {
// 	return map[string]cliCommand{
// 		"login": {
// 			name:        "login",
// 			description: "sets the current user in the config",
// 		},
// 		"register": {
// 			name:        "register",
// 			description: "adds a new user to the database",
// 		},
// 		"users": {
// 			name:        "users",
// 			description: "lists all the users in the database",
// 		},
// 	}
// }
