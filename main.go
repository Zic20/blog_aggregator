package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
	"github.com/zic20/blog_aggregator/internal"
	"github.com/zic20/blog_aggregator/internal/config"
	"github.com/zic20/blog_aggregator/internal/database"
)

func main() {

	cfg, err := config.Read()
	if err != nil {
		fmt.Printf("Error loading config: %s", err)
		os.Exit(0)
	}

	db, err := sql.Open("postgres", cfg.DBUrl)
	if err != nil {
		fmt.Printf("Error loading databas: %s", err)
		os.Exit(1)
	}

	dbQueries := database.New(db)
	state := internal.State{
		DB:     dbQueries,
		Config: &cfg,
	}

	commands := internal.Commands{
		Commands: make(map[string]func(*internal.State, internal.Command) error),
	}
	commands.Register("login", internal.HandlerLogin)
	commands.Register("register", internal.HandlerRegister)
	commands.Register("reset", internal.HandlerDelete)
	commands.Register("users", internal.HandlerGetUsers)
	commands.Register("agg", internal.HandlerAGG)
	commands.Register("addfeed", internal.HandlerAddFeed)
	commands.Register("feeds", internal.HandlerFeeds)
	commands.Register("follow", internal.HandlerFollow)

	args := os.Args

	if len(args) < 2 {
		fmt.Println("One command was expected")
		os.Exit(1)
	}
	command := internal.Command{
		Name: args[1],
		Args: args[2:],
	}

	if err = commands.Run(&state, command); err != nil {
		fmt.Println(err)
	}
}
