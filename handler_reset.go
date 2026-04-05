package main

import (
	"context"
	"fmt"
)

func handlerReset(s *state, cmd command) error {
	err := s.db.ResetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("Error resetting feeds table: %v\n", err)
	}
	err = s.db.ResetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("Error resetting users table: %v\n", err)
	}
	fmt.Println("Database reset successfully.")
	return nil
}
