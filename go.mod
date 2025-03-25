module github.com/wittyCode/blog-aggregator

go 1.24.0

require internal/config v0.0.1

require internal/rss v0.0.1

require (
	github.com/google/uuid v1.6.0
	github.com/lib/pq v1.10.9
)

replace internal/config => ./internal/config/

replace internal/rss => ./internal/rss/
