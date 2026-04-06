package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/brentjolicoeur/gator/internal/database"
	"github.com/google/uuid"
)

func handlerFollow(s *state, cmd command) error {
	user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		return fmt.Errorf("Error retrieving user: %v\n", err)
	}
	if len(cmd.args) == 0 {
		return errors.New("command syntax: follow <feed_URL>")
	}
	feedURL := cmd.args[0]

	feed_id, err := s.db.GetFeedByURL(context.Background(), feedURL)
	if err != nil {
		return fmt.Errorf("Error retrieving feed: %v\n", err)
	}

	params := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed_id,
	}
	feedFollow, err := s.db.CreateFeedFollow(context.Background(), params)
	if err != nil {
		return fmt.Errorf("Error creating feed_follow record: %v\n", err)
	}
	fmt.Printf("%s is now following %s\n", feedFollow.UserName, feedFollow.FeedName)

	return nil
}
