package rssapi

import "github.com/shhj1998/rss-search-api/rsserver"

// Server instance that has a db connection.
type Server struct {
	DB *rsserver.DB
}
