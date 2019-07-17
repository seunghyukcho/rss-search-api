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
	itemSchema = `CREATE TABLE Item (
					item_id INT NOT NULL AUTO_INCREMENT,
					guid VARCHAR(30) NOT NULL,
					title VARCHAR(255),
					link VARCHAR(255),
					description TEXT,
					pub_date DATETIME,
					creator VARCHAR(255),
					CONSTRAINT PRIMARY KEY (item_id),
					CONSTRAINT INDEX (item_id),
					CONSTRAINT UNIQUE (guid)`

	enclosureSchema = `CREATE TABLE Enclosure (
						item INT NOT NULL,
						url VARCHAR(255) NOT NULL,
						length VARCHAR(30),
						type VARCHAR(30),
						CONSTRAINT PRIMARY KEY (item, url),
						CONSTRAINT INDEX (item),
						CONSTRAINT FOREIGN KEY (item) REFERENCES Item(item_id)`

	channelSchema = `CREATE TABLE Channel (
						channel_id INT NOT NULL AUTO_INCREMENT,
						title VARCHAR(255),
						description VARCHAR(255),
						site_link VARCHAR(255),
						rss_link VARCHAR(255) NOT NULL,
						CONSTRAINT PRIMARY KEY (channel_id)
						CONSTRAINT UNIQUE (rss_link)`

	publishSchema = `CREATE TABLE Publish (
						item INT NOT NULL,
						channel INT NOT NULL,
						CONSTRAINT PRIMARY KEY (item, channel),
						CONSTRAINT INDEX (channel),
						CONSTRAINT FOREIGN KEY (item) REFERENCES Item(item_id),
						CONSTRAINT FOREIGN KEY (channel) REFERENCES Channel(channel_id)`
)
