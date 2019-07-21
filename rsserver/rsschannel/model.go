package rsschannel

import (
	"database/sql"
	"github.com/mmcdole/gofeed"
)

// Table contains the db connection to rsschannel Table.
type Table struct {
	Connection *sql.DB
}

// Schema shows the rsschannel table schema that is really used.
type Schema struct {
	ID      int    `json:"id,omitempty"`
	RSSLink string `json:"rsslink,omitempty"`
	gofeed.Feed
}
