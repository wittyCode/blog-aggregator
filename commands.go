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
