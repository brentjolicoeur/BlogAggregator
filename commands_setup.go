package main

import (
	"errors"

	"github.com/brentjolicoeur/gator/internal/config"
)

type state struct {
	config *config.Config
}

type command struct {
	name string
	args []string
}

type commands struct {
	cmds map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	command, ok := c.cmds[cmd.name]
	if !ok {
		return errors.New("Command not found")
	} else {
		return command(s, cmd)
	}
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.cmds[name] = f
}
