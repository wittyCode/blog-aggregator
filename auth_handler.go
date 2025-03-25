package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/wittyCode/blog-aggregator/internal/database"
)

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		user, err := s.db.GetUser(context.Background(), s.config.CurrentUserName)
		if err != nil {
			log.Fatal(err)
		}

		return handler(s, cmd, user)
	}
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		log.Fatal("login function expects username as argument")
	}

	userName := cmd.args[0]
	ctx := context.Background()
	_, err := s.db.GetUser(ctx, userName)
	if err != nil {
		msg := fmt.Sprintf("username %s does not exist", userName)
		log.Fatal(msg)
	}

	err = s.config.SetUser(userName)
	if err == nil {
		fmt.Println("Username has been set to ", userName)
	}

	return err
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		log.Fatal("login function expects username as argument")
	}

	userName := cmd.args[0]
	ctx := context.Background()
	now := sql.NullTime{Time: time.Now(), Valid: true}

	params := database.CreateUserParams{ID: uuid.New(), CreatedAt: now, UpdatedAt: now, Name: userName}
	_, err := s.db.CreateUser(ctx, params)
	if err != nil {
		log.Fatal(err)
	}

	s.config.SetUser(userName)
	fmt.Println("User created!")
	fmt.Println(params)

	return err
}
