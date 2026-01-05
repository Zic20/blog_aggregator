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
	"github.com/zic20/blog_aggregator/internal/rss"
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
		return err
	}

	s.Config.SetUser(user.Name)
	fmt.Println("User created successfully")
	fmt.Println(user)

	return nil
}

func HandlerGetUsers(s *State, cmd Command) error {
	users, err := s.DB.GetUsers(context.Background())
	if err != nil {
		return err
	}

	for _, user := range users {
		if user.Name == s.Config.CurrentUserName {
			fmt.Printf("* %s (current)\n", user.Name)
		} else {
			fmt.Printf("* %s\n", user.Name)
		}
	}
	return nil
}

func HandlerDelete(s *State, _ Command) error {
	if err := s.DB.DeleteAllUsers(context.Background()); err != nil {
		fmt.Printf("Reset failed: %s", err)
		os.Exit(1)
		return err
	}

	fmt.Print("Reset completed successfully")
	return nil
}

func HandlerAGG(s *State, _ Command) error {
	feed, err := rss.FetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return err
	}

	fmt.Print(*feed)

	return nil
}

func HandlerAddFeed(s *State, cmd Command) error {
	if len(cmd.Args) < 2 {
		fmt.Println("addfeed expects two arguments")
		os.Exit(1)
		return errors.New("addfeed expects two arguments")
	}

	user, err := s.DB.GetUserByName(context.Background(), s.Config.CurrentUserName)
	if err != nil {
		return fmt.Errorf("Error fetching current user: %s", err)
	}

	data := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.Args[0],
		Url:       cmd.Args[1],
		UserID:    user.ID,
	}

	feed, err := s.DB.CreateFeed(context.Background(), data)
	if err != nil {
		return err
	}

	fmt.Println(feed)
	return nil
}

func HandlerFeeds(s *State, _ Command) error {
	feeds, err := s.DB.ListFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("Error fetching feeds: %s", err)
	}
	fmt.Println(feeds)
	return nil
}
