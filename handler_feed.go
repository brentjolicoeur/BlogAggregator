package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/brentjolicoeur/gator/internal/database"
	"github.com/google/uuid"
)

func handlerAddFeed(s *state, cmd command) error {
	if len(cmd.args) < 2 {
		return errors.New("command syntax: addfeed <Name> <url>")
	}
	name := cmd.args[0]
	feedURL := cmd.args[1]
	user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		return fmt.Errorf("Error retrieving user: %v\n", err)
	}
	feedParams := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      name,
		Url:       feedURL,
		UserID:    user.ID,
	}
	feed, err := s.db.CreateFeed(context.Background(), feedParams)
	if err != nil {
		return fmt.Errorf("Error creating feed: %v\n", err)
	}
	fmt.Println("Feed created successfully")
	fmt.Printf("%+v\n", feed)
	return nil
}

func handlerListFeeds(s *state, cmd command) error {
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("Error retrieving feeds: %v\n", err)
	}

	fmt.Printf("%v feeds found in database.\n", len(feeds))
	for i, feed := range feeds {
		feed_user, err := s.db.GetUsernameById(context.Background(), feed.UserID)
		if err != nil {
			return fmt.Errorf("Error retrieving username: %v\n", err)
		}
		fmt.Printf("Feed %v name: %v\n", i+1, feed.Name)
		fmt.Printf("Feed %v url: %v\n", i+1, feed.Url)
		fmt.Printf("Feed %v created by: %v\n", i+1, feed_user)
	}
	return nil
}
