package internal

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/zic20/blog_aggregator/internal/config"
	"github.com/zic20/blog_aggregator/internal/database"
)

type State struct {
	DB     *database.Queries
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

	username := cmd.Args[0]

	_, err := s.DB.GetUserByName(context.Background(), username)
	if err != nil {
		fmt.Println("invalid user")
		os.Exit(1)
		return err
	}

	if err := s.Config.SetUser(username); err != nil {
		return err
	}
	fmt.Printf("User %s has been set", username)
	return nil
}

func HandlerRegister(s *State, cmd Command) error {
	if len(cmd.Args) == 0 {
		fmt.Println("username is required")
		os.Exit(1)
		return errors.New("register handler expects one argument")
	}

	data := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.Args[0],
	}
	user, err := s.DB.CreateUser(context.Background(), data)
	if err != nil {
		fmt.Printf("Error creating user: %s", err)
		os.Exit(1)
	}

	s.Config.SetUser(user.Name)
	fmt.Println("User created successfully")
	fmt.Println(user)

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
