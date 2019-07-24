package rssitem

import (
	"database/sql"
	"github.com/mmcdole/gofeed"
)

// Table contains the db connection to rssitem Table.
type Table struct {
	Connection *sql.DB
}

// Schema shows the rssitem table schema that is really used.
type Schema struct {
	gofeed.Item
}
