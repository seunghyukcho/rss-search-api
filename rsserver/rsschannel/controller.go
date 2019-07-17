package rsschannel

import (
	"github.com/mmcdole/gofeed"
	"github.com/shhj1998/rss-search-api/rsserver/rssitem"
)

func (table *Table) GetChannelLinks(channels *[]*Schema) (err error) {
	channelRows, err := table.Connection.Query(`SELECT channel_id, rss_link FROM Channel`)

	if err == nil {
		for channelRows.Next() {
			var link Schema
			if err := channelRows.Scan(&link.Id, &link.RSSLink); err != nil {
				return nil
			}

			*channels = append(*channels, &link)
		}
	} else {
		return err
	}

	return nil
}

func (table *Table) GetChannels(channels *[]*Schema) (err error) {
	channelRows, err := table.Connection.Query(`SELECT title, description, site_link FROM Channel`)

	if err == nil {
		if err := Fetch(channelRows, channels); err != nil {
			return err
		}
	} else {
		return err
	}

	return nil
}

func (table *Table) getChannelWithItems(channel *Schema, count int) (err error) {
	var items []*rssitem.Schema
	itemRows, err := table.Connection.Query(`SELECT I.guid, I.title, I.link, I.description, I.pub_date, I.creator, e.url, e.length, e.type 
																			FROM Channel c JOIN Publish p ON c.channel_id=p.channel JOIN Item I ON p.item=I.guid LEFT JOIN Enclosure e ON I.guid=e.item 
																			WHERE channel_id=? ORDER BY I.pub_date DESC LIMIT 0,?`, channel.Id, count)
	if err != nil {
		return err
	} else if err = rssitem.Fetch(itemRows, &items); err != nil {
		return err
	} else {
		for idx := range items {
			channel.Items = append(channel.Items, &items[idx].Item)
		}
	}

	return nil
}

func (table *Table) GetChannelsWithItems(channels *[]*Schema, count int) (err error) {
	channelRows, err := table.Connection.Query(`SELECT channel_id, title, description, site_link, rss_link FROM Channel`)

	if err != nil {
		return err
	} else {
		if err = Fetch(channelRows, channels); err != nil {
			return err
		} else {
			for idx := range *channels {
				if err = table.getChannelWithItems((*channels)[idx], count); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (table *Table) CreateChannel(rssLink string) (err error) {
	fp := gofeed.NewParser()
	if _, err = fp.ParseURL(rssLink); err == nil {
		if _, err = table.Connection.Exec(`INSERT INTO Channel(rss_link) VALUE(?)`, rssLink); err != nil {
			return err
		}
	} else {
		return err
	}

	return nil
}
