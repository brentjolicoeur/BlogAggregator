package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/brentjolicoeur/gator/internal/database"
	"github.com/google/uuid"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("Login command is missing username.\n")
	}
	name := cmd.args[0]

	_, err := s.db.GetUser(context.Background(), name)
	if err != nil {
		return errors.New("User doesn't exist")
	}

	err = s.cfg.SetUser(name)
	if err != nil {
		return err
	}
	fmt.Printf("%s has been set as user.\n", s.cfg.CurrentUserName)

	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("Register command is missing username.\n")
	}
	name := cmd.args[0]

	_, err := s.db.GetUser(context.Background(), name)
	if err == nil {
		return errors.New("User already exists.")
	}

	params := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      name,
	}

	user, err := s.db.CreateUser(context.Background(), params)
	if err != nil {
		return errors.New("error creating user")
	}
	err = s.cfg.SetUser(user.Name)
	if err != nil {
		return errors.New("error setting username")
	}
	fmt.Println("User successfully created.")

	return nil
}

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

func handlerUsers(s *state, cmd command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("Error retrieving users: %v\n", err)
	}
	for _, user := range users {
		if s.cfg.CurrentUserName == user.Name {
			fmt.Printf(" * %v (current)\n", user.Name)
		} else {
			fmt.Printf(" * %v\n", user.Name)
		}
	}
	return nil
}

func handlerAgg(s *state, cmd command) error {
	feedURL := "https://www.wagslane.dev/index.xml"

	feed, err := fetchFeed(context.Background(), feedURL)
	if err != nil {
		return fmt.Errorf("Error fetching feed: %v\n", err)
	}
	fmt.Println("Feed fectched successfully.")
	fmt.Printf("%+v\n", feed)
	return nil
}

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

func handlerFeeds(s *state, cmd command) error {
	feeds, err := s.db.GetFeedsInfo(context.Background())
	if err != nil {
		return fmt.Errorf("Error retrieving feeds: %v\n", err)
	}
	fmt.Printf("%v feeds found in database.\n", len(feeds))
	for i, feed := range feeds {
		fmt.Printf("Feed %v name: %v\n", i+1, feed.Name)
		fmt.Printf("Feed %v url: %v\n", i+1, feed.Url)
		fmt.Printf("Feed %v created by: %v\n", i+1, feed.User)
	}
	return nil
}
