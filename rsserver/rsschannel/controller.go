// Package rsschannel provides functionality to act with
// rsserver's channel Table.
package rsschannel

import (
	"database/sql"
	"github.com/mmcdole/gofeed"
	"github.com/shhj1998/rss-search-api/rsserver/handle"
	"github.com/shhj1998/rss-search-api/rsserver/rssitem"
)

// SelectLinks fetch the rss links of channels with its id
// and store it in channels array.
func (table *Table) SelectLinks(channels *[]*Schema) (err error) {
	var tx *sql.Tx
	if tx, err = table.Connection.Begin(); err != nil {
		return err
	}

	channelRows, err := tx.Query(`SELECT channel_id, rss_link FROM Channel`)

	if err == nil {
		for channelRows.Next() {
			var link Schema
			if err := channelRows.Scan(&link.ID, &link.RSSLink); err != nil {
				return handle.Transaction(tx, err)
			}

			*channels = append(*channels, &link)
		}
	} else {
		return handle.Transaction(tx, err)
	}

	return tx.Commit()
}

// Select fetch the information of channels in the database
// and store it in channels array. You can see type Schema
// in model.go to check the features.
func (table *Table) Select(channels *[]*Schema) (err error) {
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

func (table *Table) selectChannelWithItems(channel *Schema, count int) (err error) {
	var tx *sql.Tx
	var items []*rssitem.Schema

	if tx, err = table.Connection.Begin(); err != nil {
		return err
	}

	itemRows, err := tx.Query(`SELECT I.guid, I.title, I.link, I.description, I.pub_date, I.creator, e.url, e.length, e.type 
																			FROM Channel c JOIN Publish p ON c.channel_id=p.channel JOIN Item I ON p.item=I.item_id LEFT JOIN Enclosure e ON I.item_id=e.item 
																			WHERE channel_id=? ORDER BY I.pub_date DESC LIMIT 0,?`, channel.ID, count)
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

func containValue(ids *[]int, find int) bool {
	if ids == nil {
		return false
	}

	for _, val := range *ids {
		if val == find {
			return true
		}
	}

	return false
}

// SelectWithItems fetch channels with its items. You can
// limit the channels and items number by the parameter ids and count.
func (table *Table) SelectWithItems(channels *[]*Schema, ids *[]int, count int) (err error) {
	var tx *sql.Tx
	var totalChannels []*Schema
	var channelRows *sql.Rows
	if tx, err = table.Connection.Begin(); err != nil {
		return err
	}

	if channelRows, err = tx.Query(`SELECT channel_id, title, description, site_link, rss_link FROM Channel`); err != nil {
		return handle.Transaction(tx, err)
	}

	if err = Fetch(channelRows, &totalChannels); err != nil {
		return handle.Transaction(tx, err)
	}
	for _, channel := range totalChannels {
		if ids != nil && !containValue(ids, channel.ID) {
			continue
		}

		if err = table.selectChannelWithItems(channel, count); err != nil {
			return handle.Transaction(tx, err)
		}

		*channels = append(*channels, channel)
	}

	return tx.Commit()
}

// Create makes a new channel instance with the rsslink.
// It is stored in the configured database.
func (table *Table) Create(rssLink string) (err error) {
	var tx *sql.Tx
	if tx, err = table.Connection.Begin(); err != nil {
		return err
	}

	fp := gofeed.NewParser()
	if _, err = fp.ParseURL(rssLink); err == nil {
		if _, err = tx.Exec(`INSERT INTO Channel (rss_link) VALUE(?)`, rssLink); err != nil {
			return handle.Transaction(tx, err)
		}
	} else {
		return handle.Transaction(tx, err)
	}

	return tx.Commit()
}
