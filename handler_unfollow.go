package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/brentjolicoeur/gator/internal/database"
)

func handlerUnfollowFeed(s *state, cmd command, user database.User) error {
	if len(cmd.args) == 0 {
		return errors.New("command syntax: unfollow <feed_URL>")
	}
	feedURL := cmd.args[0]

	feedID, err := s.db.GetFeedByURL(context.Background(), feedURL)
	if err != nil {
		return fmt.Errorf("Error retrieving feed: %v\n", err)
	}

	params := database.UnfollowFeedParams{
		UserID: user.ID,
		FeedID: feedID,
	}

	err = s.db.UnfollowFeed(context.Background(), params)
	if err != nil {
		return fmt.Errorf("Error unfollowing feed: %v\n", err)
	}
	fmt.Println("Feed successfully unfollowed.")
	return nil
}
