package rsserver

import (
	"database/sql"
	"fmt"
	"github.com/mmcdole/gofeed"
	"github.com/shhj1998/rss-search-api/rsserver/rsschannel"
	"github.com/shhj1998/rss-search-api/rsserver/rssitem"
	"time"
)

func (db *DB) Open(name, address, id, password string) (err error) {
	dbInfo := fmt.Sprintf("%s:%s@tcp(%s)/%s", id, password, address, name)
	if conn, err := sql.Open("mysql", dbInfo); err != nil {
		return err
	} else {
		db.connection = conn
		db.Name = name
		db.ItemTable = &rssitem.Table{conn}
		db.ChannelTable = &rsschannel.Table{conn}

		return nil
	}
}

func (db *DB) Create() {

}

func (db *DB) Close() (err error) {
	err = db.connection.Close()
	return err
}

func (db *DB) Update() (err error) {
	parser := gofeed.NewParser()
	conn := db.connection

	updateChannel, _ := conn.Prepare(`UPDATE Channel SET title=?, description=?, site_link=? WHERE channel_id=?`)
	insertItem, _ := conn.Prepare(`INSERT INTO Item (guid, title, description, link, pub_date, creator) VALUES (?, ?, ?, ?, ?, ?)
											ON DUPLICATE KEY UPDATE title=?, description=?, link=?, pub_date=?, creator=?`)
	insertEnclosure, _ := conn.Prepare(`INSERT IGNORE INTO Enclosure (item, url, length, type) VALUES (?, ?, ?, ?)`)
	insertPublish, _ := conn.Prepare(`INSERT IGNORE INTO Publish (item, channel) VALUES (?, ?)`)

	if channelRows, err := db.connection.Query(`SELECT channel_id, rss_link FROM Channel`); err == nil {
		for channelRows.Next() {
			var channelID int
			var rssLink string

			if err := channelRows.Scan(&channelID, &rssLink); err != nil {
				return err
			}

			channel, _ := parser.ParseURL(rssLink)

			if _, err := updateChannel.Exec(channel.Title, channel.Description, channel.Link, channelID); err != nil {
				return err
			}

			for _, item := range channel.Items {
				var author string
				publishedDate, _ := time.Parse(time.RFC3339, item.PublishedParsed.Format(time.RFC3339))

				if item.Author == nil {
					author = ""
				} else {
					author = item.Author.Name
				}

				if _, err := insertItem.Exec(item.GUID, item.Title, item.Description, item.Link, publishedDate, author, item.Title, item.Description, item.Link, publishedDate, author); err != nil {
					return err
				}

				if _, err := insertPublish.Exec(item.GUID, channelID); err != nil {
					return err
				}

				for _, enclosure := range item.Enclosures {
					if _, err := insertEnclosure.Exec(item.GUID, enclosure.URL, enclosure.Length, enclosure.Type); err != nil {
						return err
					}
				}
			}
		}
		return nil
	} else {
		return err
	}
}
