package main

import (
	"errors"
	"fmt"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("Login command is missing username.\n")
	}
	name := cmd.args[0]

	err := s.cfg.SetUser(name)
	if err != nil {
		return err
	}
	fmt.Printf("%s has been set as user.\n", s.cfg.CurrentUserName)

	return nil
}
