package main

import (
	"errors"

	"github.com/brentjolicoeur/gator/internal/config"
	"github.com/brentjolicoeur/gator/internal/database"
)

type state struct {
	db  *database.Queries
	cfg *config.Config
}

type command struct {
	name string
	args []string
}

type commands struct {
	registeredCommands map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	command, ok := c.registeredCommands[cmd.name]
	if !ok {
		return errors.New("Command not found\n")
	} else {
		return command(s, cmd)
	}
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.registeredCommands[name] = f
}
