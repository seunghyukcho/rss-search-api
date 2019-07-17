package rsschannel

import (
	"database/sql"
	"github.com/mmcdole/gofeed"
)

type Table struct {
	Connection *sql.DB
}

type Schema struct {
	Id      int    `json:"id,omitempty"`
	RSSLink string `json:"rsslink,omitempty"`
	gofeed.Feed
}
