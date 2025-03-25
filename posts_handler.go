package main

import (
	"context"
	"database/sql"
	"fmt"
	"internal/rss"
	"log"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/wittyCode/blog-agggregator/internal/database"
)

func handlerBrowse(s *state, cmd command, user database.User) error {
	limit := 2
	var err error
	if len(cmd.args) > 0 {
		limit, err = strconv.Atoi(cmd.args[0])
		if err != nil {
			log.Fatal("could not map given limmit parameter to integer, you need to provide a valid number")
		}
	}

	ctx := context.Background()
	params := database.GetPostsForUserParams{ID: user.ID, Limit: int32(limit)}
	posts, err := s.db.GetPostsForUser(ctx, params)
	if err != nil {
		log.Fatal(err)
	}

	printPosts(posts)

	return nil
}

func printPosts(posts []database.Post) {
	for _, post := range posts {
		fmt.Printf("%s (%s)\n", post.Title, post.Url)
	}
}

func persistPost(s *state, ctx context.Context, rssItem rss.RSSItem, feedId uuid.UUID) error {
	now := sql.NullTime{Time: time.Now(), Valid: true}
	publishedDate, err := parsePublishedDate(rssItem.PubDate)
	var pubDate sql.NullTime
	if err != nil {
		pubDate = now
	} else {
		pubDate = sql.NullTime{Time: publishedDate, Valid: true}
	}

	params :=
		database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   now,
			UpdatedAt:   now,
			Title:       rssItem.Title,
			Url:         rssItem.Link,
			Description: sql.NullString{String: rssItem.Description, Valid: true},
			PublishedAt: pubDate,
			FeedID:      feedId,
		}

	_, err = s.db.CreatePost(ctx, params)
	if err != nil {
		// Check if it's a unique violation
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" {
				// This is a duplicate URL error
				// You can check pqErr.Constraint to see which constraint was violated
				return nil // or handle as you see fit
			}
		}
		// Handle other errors
		log.Printf("Error creating post: %v", err)
	}
	return nil
}
