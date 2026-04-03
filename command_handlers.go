package main

import (
	"errors"
	"fmt"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("Login command is missing username.")
	}

	err := s.config.SetUser(cmd.args[0])
	if err != nil {
		return err
	}
	fmt.Printf("%s has been set as user.\n", s.config.CurrentUserName)

	return nil
}
