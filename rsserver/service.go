package rsserver

import (
	"database/sql"
	"fmt"
	"github.com/mmcdole/gofeed"
	"github.com/shhj1998/rss-search-api/rsserver/handle"
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

func (db *DB) Create() (err error) {
	tx, err := db.connection.Begin()
	if err != nil {
		return err
	}

	if _, err := tx.Exec(itemSchema); err != nil {
		return handle.Transaction(tx, err)
	}

	if _, err := tx.Exec(enclosureSchema); err != nil {
		return handle.Transaction(tx, err)
	}

	if _, err := tx.Exec(channelSchema); err != nil {
		return handle.Transaction(tx, err)
	}

	if _, err := tx.Exec(publishSchema); err != nil {
		return handle.Transaction(tx, err)
	}

	return tx.Commit()
}

func (db *DB) Close() (err error) {
	err = db.connection.Close()
	return err
}

func (db *DB) Update() (err error) {
	var tx *sql.Tx
	if tx, err = db.connection.Begin(); err != nil {
		return err
	}

	parser := gofeed.NewParser()

	updateChannel, _ := db.connection.Prepare(`UPDATE Channel SET title=?, description=?, site_link=? WHERE channel_id=?`)
	insertItem, _ := db.connection.Prepare(`INSERT INTO Item (guid, title, description, link, pub_date, creator) VALUES (?, ?, ?, ?, ?, ?)
											ON DUPLICATE KEY UPDATE title=?, description=?, link=?, pub_date=?, creator=?`)
	insertEnclosure, _ := db.connection.Prepare(`INSERT IGNORE INTO Enclosure (item, url, length, type) VALUES (?, ?, ?, ?)`)
	insertPublish, _ := db.connection.Prepare(`INSERT IGNORE INTO Publish (item, channel) VALUES (?, ?)`)

	if channelRows, err := tx.Query(`SELECT channel_id, rss_link FROM Channel`); err == nil {
		for channelRows.Next() {
			var channelID int
			var rssLink string

			if err := channelRows.Scan(&channelID, &rssLink); err != nil {
				return handle.Transaction(tx, err)
			}

			channel, _ := parser.ParseURL(rssLink)

			if _, err := updateChannel.Exec(channel.Title, channel.Description, channel.Link, channelID); err != nil {
				return handle.Transaction(tx, err)
			}

			for _, item := range channel.Items {
				var author string
				publishedDate, _ := time.Parse(time.RFC3339, item.PublishedParsed.Format(time.RFC3339))

				if item.Author == nil {
					author = ""
				} else {
					author = item.Author.Name
				}

				var result sql.Result
				if result, err = insertItem.Exec(item.GUID, item.Title, item.Description, item.Link, publishedDate, author, item.Title, item.Description, item.Link, publishedDate, author); err != nil {
					return handle.Transaction(tx, err)
				}

				success, _ := result.RowsAffected()
				if success == 1 {
					id, _ := result.LastInsertId()

					if _, err := insertPublish.Exec(id, channelID); err != nil {
						return handle.Transaction(tx, err)
					}

					for _, enclosure := range item.Enclosures {
						if _, err := insertEnclosure.Exec(id, enclosure.URL, enclosure.Length, enclosure.Type); err != nil {
							return handle.Transaction(tx, err)
						}
					}
				}
			}
		}

		return tx.Commit()
	} else {
		return handle.Transaction(tx, err)
	}
}
