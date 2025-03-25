# Gator - RSS Feed aggregator based on Golang

This is a nice little golang training project, implementing a rss feed aggregator. 

## Dependencies
You need to install PostgreSQL and Go locally to run this gator

## Installation
```bash
go get github.com/wittyCode/blog-aggregator
go install github.com/wittyCode/blog-aggregator
```

## Local Setup
To run Gator, you need to place a configuration file in your home directory called .gatorconfig.json

```json
{
    "db_url": "<YOUR_POSTGRESQL_CONNECTION_STRING_HERE>"
}
```
hint: you might need to have the suffix ?sslmode=disable

Additionally you need to run the schema migrations contained in sql/schema with goose, to make sure the necessary tables are available in your DB at runtime

```bash
goose "<YOUR_POSTGRESQL_CONNECTION_STRING_HERE>" up
```

## Usage
You have several commands available to use the blog-aggregator

### User Management

```bash
blog-aggregator register <user_name>
```

Register a new user in the database with the given name.

```bash
blog-aggregator login <user_name>
```

Log in with given name, user must exist

```bash
blog-aggregator users
```

Print the registered users in the DB

```bash
blog-aggregator reset
```

Delete all data. This will delete all users and by cascading deletes all feeds and posts as well

### RSS Feed Management

```bash
blog-aggregator addfeed <feed_name> <url>
```

Add a new RSS Feed for the given name and url to the database. The logged in user will automatically follow the newly created feed.

```bash
blog-aggregator feeds
```

List all available feeds

```bash
blog-aggregator follow <feed_url>
```

Follow the feed for the given url with the currently logged in user. Feed must exist (if it doesn't, create it first with the addfeed command)

```bash
blog-aggregator following
```

Print a list of Feeds the currently logged in user is following

```bash
blog-aggregator unfollow <feed_url>
```

No longer follow the feed for the given url with the currently logged in user.

```bash
blog-aggregator agg <time_pattern>
```

Activtate the aggregator function. It will automatically trigger feed updates in the given time period, e.g. 10s, 10m, 1h.
Time Pattern must be valid, don't set it too low to not spam the corresponding RSS Feeds with requests. Quit execution with Ctrl-C.

```bash
blog-aggregator browse [limit]
```

Browse available posts across RSS Feeds the currently logged in user follows. Most recent posts are shown first, the optional limit parameter defines how many posts should be shown at maximum. If no limit is given, the default is "2".
