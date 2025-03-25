package main

import (
	"context"
	"database/sql"
	"fmt"
	"internal/rss"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/wittyCode/blog-agggregator/internal/database"
)

func handlerAgg(s *state, cmd command) error {
	ctx := context.Background()
	rssFeed, err := rss.FetchFeed(ctx, "https://www.wagslane.dev/index.xml")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(rssFeed)
	return nil
}

func handlerAddFeed(s *state, cmd command) error {
	if len(cmd.args) < 2 {
		log.Fatal("did not provide enough arguments to addFeed, need feed name and feed url")
	}

	ctx := context.Background()
	currentUser := getCurrentUser(s, ctx)

	now := sql.NullTime{Time: time.Now(), Valid: true}
	params := database.CreateFeedParams{ID: uuid.New(), CreatedAt: now, UpdatedAt: now, Name: cmd.args[0], Url: cmd.args[1], UserID: currentUser.ID}

	newFeed, err := s.db.CreateFeed(ctx, params)

	if err != nil {
		log.Fatal(err)
	}

	// Create a feed follow directly instead of calling handlerFollow
	followParams := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: now,
		UpdatedAt: now,
		UserID:    currentUser.ID,
		FeedID:    newFeed.ID,
	}

	feedFollow, err := s.db.CreateFeedFollow(ctx, followParams)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("You are now following '%s'\n", feedFollow.FeedName)

	return nil
}

func handlerFeeds(s *state, cmd command) error {
	var feeds []database.Feed
	ctx := context.Background()
	feeds, err := s.db.GetFeeds(ctx)
	if err != nil {
		log.Fatal(err)
	}

	userIds := make([]uuid.UUID, len(feeds))
	for i := range feeds {
		userIds[i] = feeds[i].UserID
	}

	users, err := s.db.GetUsersByIds(ctx, userIds)
	if err != nil {
		log.Fatal(err)
	}

	usersById := make(map[uuid.UUID]database.User)
	for _, user := range users {
		usersById[user.ID] = user
	}

	printFeeds(feeds, usersById)

	return nil
}

func printFeeds(feeds []database.Feed, usersById map[uuid.UUID]database.User) {
	for _, feed := range feeds {
		user := usersById[feed.UserID]
		fmt.Printf("%s (%s); created by %s\n", feed.Name, feed.Url, user.Name)
	}
}

func handlerFollow(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		log.Fatal("follow command expects url of rss feed to follow as argument")
	}

	ctx := context.Background()
	currentUser := getCurrentUser(s, ctx)

	url := cmd.args[0]
	feed, err := s.db.GetFeedByUrl(ctx, url)
	if err != nil {
		msg := fmt.Sprintf("feed %s does not exist", url)
		log.Fatal(msg)
	}

	now := sql.NullTime{Time: time.Now(), Valid: true}
	feedFollow := database.CreateFeedFollowParams{ID: uuid.New(), CreatedAt: now, UpdatedAt: now, UserID: currentUser.ID, FeedID: feed.ID}

	insertedFeedFollow, err := s.db.CreateFeedFollow(ctx, feedFollow)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s is now following the feed %s\n", insertedFeedFollow.UserName, insertedFeedFollow.FeedName)

	return nil
}

func getCurrentUser(s *state, ctx context.Context) database.User {
	currentUser, err := s.db.GetUser(ctx, s.config.CurrentUserName)

	if err != nil {
		msg := fmt.Sprintf("username %s does not exist", s.config.CurrentUserName)
		log.Fatal(msg)
	}

	return currentUser
}

func handlerFollowing(s *state, cmd command) error {
	ctx := context.Background()
	currentUser := getCurrentUser(s, ctx)

	feedsOfCurrentUser, err := s.db.GetFeedFollowsForUser(ctx, currentUser.ID)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s is following these feeds:\n", currentUser.Name)
	for _, feed := range feedsOfCurrentUser {
		fmt.Println(feed.FeedName)
	}

	return nil
}
