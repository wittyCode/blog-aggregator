package main

import (
	"context"
	"fmt"
	"log"

	"github.com/wittyCode/blog-aggregator/internal/database"
)

func handlerReset(s *state, cmd command) error {
	ctx := context.Background()
	err := s.db.ResetUsers(ctx)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("sucessfully reset users table in DB")
	return err
}

func handlerUsers(s *state, cmd command) error {
	ctx := context.Background()
	users, err := s.db.GetUsers(ctx)
	if err != nil {
		log.Fatal(err)
	}

	printUsers(users, s.config.CurrentUserName)

	return err
}

func printUsers(users []database.User, currentUser string) {
	for _, user := range users {
		if user.Name != currentUser {
			fmt.Printf("* %s\n", user.Name)
		} else {
			fmt.Printf("* %s (current)\n", user.Name)
		}
	}
}
