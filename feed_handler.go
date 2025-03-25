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
	currentUser, err := s.db.GetUser(ctx, s.config.CurrentUserName)

	if err != nil {
		msg := fmt.Sprintf("username %s does not exist", s.config.CurrentUserName)
		log.Fatal(msg)
	}

	now := sql.NullTime{Time: time.Now(), Valid: true}
	params := database.CreateFeedParams{ID: uuid.New(), CreatedAt: now, UpdatedAt: now, Name: cmd.args[0], Url: cmd.args[1], UserID: currentUser.ID}

	newFeed, err := s.db.CreateFeed(ctx, params)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(newFeed)

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
