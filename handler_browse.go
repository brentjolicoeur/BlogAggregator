package main

import (
	"context"
	"fmt"
	"strconv"

	"github.com/brentjolicoeur/gator/internal/database"
)

func handlerBrowsePosts(s *state, cmd command, user database.User) error {
	var postLimit int32

	if len(cmd.args) == 0 {
		fmt.Println("Syntax: browse <number_of_posts>  (Default is 2)")
		postLimit = 2
	} else {
		limit, err := strconv.Atoi(cmd.args[0])
		if err != nil {
			fmt.Println("Supplied arguemtn not a valid number.")
			fmt.Println("Syntax: browse <number_of_posts>")
			return err
		}
		postLimit = int32(limit)
	}

	params := database.GetPostsUserParams{
		UserID: user.ID,
		Limit:  postLimit,
	}
	posts, err := s.db.GetPostsUser(context.Background(), params)
	if err != nil {
		return fmt.Errorf("Error retrieving posts: %v\n", err)
	}
	for _, post := range posts {
		fmt.Printf("Title: %v\n", post.Title.String)
		fmt.Printf("%v\n", post.Description.String)
	}
	return nil
}
