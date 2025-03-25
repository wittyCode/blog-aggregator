package main

import (
	"database/sql"
	"internal/config"
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/wittyCode/blog-aggregator/internal/database"
)

func main() {
	configFile := config.Read()
	st := state{config: &configFile}

	db, err := sql.Open("postgres", st.config.DbUrl)
	if err != nil {
		log.Fatal("error connecting to database with connection string ", st.config.DbUrl)
	}

	dbQueries := database.New(db)
	st.db = dbQueries

	cmds := commands{make(map[string]func(*state, command) error)}
	cmds.registerCommands()

	args := os.Args
	if len(args) < 2 {
		log.Fatal("Too few arguments, we need at least 2 arguments")
	}

	cmdName := args[1]
	params := args[2:]

	cmd := command{cmdName, params}

	cmds.run(&st, cmd)
}
