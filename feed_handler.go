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
	if len(cmd.args) == 0 {
		log.Fatal("command agg needs time duration param in format 1s, 1m, 1h...")
	}

	timeBetweenRequests, err := time.ParseDuration(cmd.args[0])
	if err != nil {
		log.Fatal("invalid format of given duration param")
	}

	fmt.Printf("Collecting feeds every %v\n", timeBetweenRequests)

	ticker := time.NewTicker(timeBetweenRequests)
	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}
}

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 2 {
		log.Fatal("did not provide enough arguments to addFeed, need feed name and feed url")
	}

	ctx := context.Background()

	now := sql.NullTime{Time: time.Now(), Valid: true}
	params := database.CreateFeedParams{ID: uuid.New(), CreatedAt: now, UpdatedAt: now, Name: cmd.args[0], Url: cmd.args[1], UserID: user.ID}

	newFeed, err := s.db.CreateFeed(ctx, params)

	if err != nil {
		log.Fatal(err)
	}

	// Create a feed follow directly instead of calling handlerFollow
	followParams := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: now,
		UpdatedAt: now,
		UserID:    user.ID,
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

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) == 0 {
		log.Fatal("follow command expects url of rss feed to follow as argument")
	}

	ctx := context.Background()

	url := cmd.args[0]
	feed, err := s.db.GetFeedByUrl(ctx, url)
	if err != nil {
		msg := fmt.Sprintf("feed %s does not exist", url)
		log.Fatal(msg)
	}

	now := sql.NullTime{Time: time.Now(), Valid: true}
	feedFollow := database.CreateFeedFollowParams{ID: uuid.New(), CreatedAt: now, UpdatedAt: now, UserID: user.ID, FeedID: feed.ID}

	insertedFeedFollow, err := s.db.CreateFeedFollow(ctx, feedFollow)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s is now following the feed %s\n", insertedFeedFollow.UserName, insertedFeedFollow.FeedName)

	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {
	ctx := context.Background()

	feedsOfCurrentUser, err := s.db.GetFeedFollowsForUser(ctx, user.ID)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s is following these feeds:\n", user.Name)
	for _, feed := range feedsOfCurrentUser {
		fmt.Println(feed.FeedName)
	}

	return nil
}

func handlerUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) == 0 {
		log.Fatal("command unfollow expects url argument")
	}

	url := cmd.args[0]
	ctx := context.Background()
	feed, err := s.db.GetFeedByUrl(ctx, url)
	if err != nil {
		log.Fatal(err)
	}

	params := database.DeleteFollowParams{UserID: user.ID, FeedID: feed.ID}
	err = s.db.DeleteFollow(ctx, params)

	return nil
}

func scrapeFeeds(s *state) {
	ctx := context.Background()
	feed, err := s.db.GetNextFeedToFetch(ctx)
	if err != nil {
		fmt.Println(err)
		return
	}

	now := sql.NullTime{Time: time.Now(), Valid: true}
	params := database.MarkFeedFetchedParams{UpdatedAt: now, LastFetchedAt: now, ID: feed.ID}
	err = s.db.MarkFeedFetched(ctx, params)
	if err != nil {
		fmt.Println(err)
		return
	}

	rssFeed, err := rss.FetchFeed(ctx, feed.Url)

	for _, item := range rssFeed.Channel.Item {
		fmt.Println(item.Title)
	}
}
