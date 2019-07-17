package rsserver

import (
	"database/sql"
	"github.com/shhj1998/rss-search-api/rsserver/rsschannel"
	"github.com/shhj1998/rss-search-api/rsserver/rssitem"
)

type DB struct {
	connection   *sql.DB
	Name         string
	ItemTable    *rssitem.Table
	ChannelTable *rsschannel.Table
}
