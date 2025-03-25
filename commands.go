package main

import (
	"fmt"
	"internal/config"

	"github.com/wittyCode/blog-agggregator/internal/database"
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

func (cmds commands) register(name string, f func(*state, command) error) {
	cmds.commandMap[name] = f
}

func (cmds commands) run(s *state, cmd command) error {
	if cmdFunc, ok := cmds.commandMap[cmd.name]; ok {
		return cmdFunc(s, cmd)
	}

	return fmt.Errorf("command %s does not exist!", cmd.name)
}

func (cmds commands) registerCommands() {
	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerUsers)
	cmds.register("agg", handlerAgg)
	cmds.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	cmds.register("feeds", handlerFeeds)
	cmds.register("follow", middlewareLoggedIn(handlerFollow))
	cmds.register("following", middlewareLoggedIn(handlerFollowing))
	cmds.register("unfollow", middlewareLoggedIn(handlerUnfollow))
}
