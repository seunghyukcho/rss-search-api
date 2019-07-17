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

const (
	itemSchema = `CREATE TABLE IF NOT EXISTS Item (
					item_id INT NOT NULL AUTO_INCREMENT,
					guid VARCHAR(30) NOT NULL UNIQUE,
					title VARCHAR(255),
					link VARCHAR(255),
					description TEXT,
					pub_date DATETIME,
					creator VARCHAR(255),
					INDEX (item_id),
					PRIMARY KEY (item_id))`

	enclosureSchema = `CREATE TABLE IF NOT EXISTS Enclosure (
						item INT NOT NULL,
						url VARCHAR(255) NOT NULL UNIQUE,
						length VARCHAR(30),
						type VARCHAR(30),
						INDEX (item),
						PRIMARY KEY (item, url),
						FOREIGN KEY (item) REFERENCES Item(item_id))`

	channelSchema = `CREATE TABLE IF NOT EXISTS Channel (
						channel_id INT NOT NULL AUTO_INCREMENT,
						title VARCHAR(255) DEFAULT "",
						description VARCHAR(255) DEFAULT "",
						site_link VARCHAR(255) DEFAULT "",
						rss_link VARCHAR(255) NOT NULL UNIQUE,
						INDEX (channel_id),
						PRIMARY KEY (channel_id))`

	publishSchema = `CREATE TABLE IF NOT EXISTS Publish (
						item INT NOT NULL,
						channel INT NOT NULL,
						PRIMARY KEY (item, channel),
						FOREIGN KEY (item) REFERENCES Item(item_id),
						FOREIGN KEY (channel) REFERENCES Channel(channel_id))`
)
