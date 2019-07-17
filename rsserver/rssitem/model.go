package rssitem

import (
	"database/sql"
	"github.com/mmcdole/gofeed"
)

type Table struct {
	Connection *sql.DB
}

type Schema struct {
	gofeed.Item
}
