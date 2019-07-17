package rsschannel

import (
	"database/sql"
	"github.com/mmcdole/gofeed"
	"github.com/shhj1998/rss-search-api/rsserver/handle"
	"github.com/shhj1998/rss-search-api/rsserver/rssitem"
)

func (table *Table) GetLinks(channels *[]*Schema) (err error) {
	var tx *sql.Tx
	if tx, err = table.Connection.Begin(); err != nil {
		return err
	}

	channelRows, err := tx.Query(`SELECT channel_id, rss_link FROM Channel`)

	if err == nil {
		for channelRows.Next() {
			var link Schema
			if err := channelRows.Scan(&link.Id, &link.RSSLink); err != nil {
				return handle.Transaction(tx, err)
			}

			*channels = append(*channels, &link)
		}
	} else {
		return handle.Transaction(tx, err)
	}

	return tx.Commit()
}

func (table *Table) Get(channels *[]*Schema) (err error) {
	var tx *sql.Tx
	if tx, err = table.Connection.Begin(); err != nil {
		return err
	}

	channelRows, err := tx.Query(`SELECT channel_id, title, description, site_link, rss_link FROM Channel`)

	if err == nil {
		if err := Fetch(channelRows, channels); err != nil {
			return handle.Transaction(tx, err)
		}
	} else {
		return handle.Transaction(tx, err)
	}

	return tx.Commit()
}

func (table *Table) getChannelWithItems(channel *Schema, count int) (err error) {
	var tx *sql.Tx
	var items []*rssitem.Schema

	if tx, err = table.Connection.Begin(); err != nil {
		return err
	}

	itemRows, err := tx.Query(`SELECT I.guid, I.title, I.link, I.description, I.pub_date, I.creator, e.url, e.length, e.type 
																			FROM Channel c JOIN Publish p ON c.channel_id=p.channel JOIN Item I ON p.item=I.item_id LEFT JOIN Enclosure e ON I.guid=e.item 
																			WHERE channel_id=? ORDER BY I.pub_date DESC LIMIT 0,?`, channel.Id, count)
	if err != nil {
		return handle.Transaction(tx, err)
	} else if err = rssitem.Fetch(itemRows, &items); err != nil {
		return handle.Transaction(tx, err)
	} else {
		for idx := range items {
			channel.Items = append(channel.Items, &items[idx].Item)
		}
	}

	return tx.Commit()
}

func (table *Table) GetWithItems(channels *[]*Schema, count int) (err error) {
	var tx *sql.Tx
	if tx, err = table.Connection.Begin(); err != nil {
		return err
	}

	if channelRows, err := tx.Query(`SELECT channel_id, title, description, site_link, rss_link FROM Channel`); err != nil {
		return handle.Transaction(tx, err)
	} else {
		if err = Fetch(channelRows, channels); err != nil {
			return handle.Transaction(tx, err)
		} else {
			for idx := range *channels {
				if err = table.getChannelWithItems((*channels)[idx], count); err != nil {
					return handle.Transaction(tx, err)
				}
			}
		}
	}

	return tx.Commit()
}

func (table *Table) Create(rssLink string) (err error) {
	var tx *sql.Tx
	if tx, err = table.Connection.Begin(); err != nil {
		return err
	}

	fp := gofeed.NewParser()
	if _, err = fp.ParseURL(rssLink); err == nil {
		if _, err = tx.Exec(`INSERT INTO Channel(rss_link) VALUE(?)`, rssLink); err != nil {
			return handle.Transaction(tx, err)
		}
	} else {
		return handle.Transaction(tx, err)
	}

	return tx.Commit()
}
