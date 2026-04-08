package main

import (
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"html"
	"io"
	"net/http"
	"time"

	"github.com/brentjolicoeur/gator/internal/database"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	client := http.Client{
		Timeout: 5 * time.Second,
	}
	var feed RSSFeed

	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "gator")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if err = xml.Unmarshal(body, &feed); err != nil {
		return nil, err
	}
	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)
	for i := 0; i < len(feed.Channel.Item); i++ {
		feed.Channel.Item[i].Title = html.UnescapeString(feed.Channel.Item[i].Title)
		feed.Channel.Item[i].Description = html.UnescapeString(feed.Channel.Item[i].Description)
	}

	return &feed, nil
}

func scrapeFeeds(s *state) error {
	ctx := context.Background()

	nextFeed, err := s.db.GetNextFeedToFetch(ctx)
	if err != nil {
		return fmt.Errorf("Error finding next feed: %v\n", err)
	}
	fmt.Printf("Fetching items from : %v\n", nextFeed.Name)

	feed, err := fetchFeed(ctx, nextFeed.Url)
	if err != nil {
		return fmt.Errorf("Error fetching feed: %v\n", err)
	}

	for _, item := range feed.Channel.Item {
		postParams := database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Title:       convertToNullString(item.Title),
			Url:         item.Link,
			Description: convertToNullString(item.Description),
			PublishedAt: parsePublishedAt(item.PubDate),
			FeedID:      nextFeed.ID,
		}
		_, err := s.db.CreatePost(ctx, postParams)
		if err != nil {
			var pqErr *pq.Error
			if errors.As(err, &pqErr) && pqErr.Code == "23505" {
				continue
			} else {
				return fmt.Errorf("Error creating post: %v\n", err)
			}
		}
		fmt.Println("New post created.")
	}
	err = s.db.MarkFeedFetched(ctx, nextFeed.ID)
	if err != nil {
		return fmt.Errorf("Error marking feed as fetched: %v\n", err)
	}
	return nil
}
