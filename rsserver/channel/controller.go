package channel

import (
	"github.com/mmcdole/gofeed"
	"github.com/shhj1998/rss-search-api/rsserver"
	"github.com/shhj1998/rss-search-api/rsserver/item"
)

type Controller struct {
	Table *rsserver.DB
}

func (controller *Controller) GetChannels(channels *[]Channel) (err error) {
	channelRows, err := controller.Table.Connection.Query(`SELECT channel_id, title, description, site_link, rss_link FROM Channel`)

	if err == nil {
		if err := Fetch(channelRows, channels); err != nil {
			return err
		}
	} else {
		return err
	}

	return nil
}

func (controller *Controller) getChannelWithItems(channel *Channel, count int) (err error) {
	itemRows, err := controller.Table.Connection.Query(`SELECT I.guid, I.title, I.link, I.description, I.pub_date, I.creator, e.url, e.length, e.type 
																			FROM Channel c JOIN Publish p ON c.channel_id=p.channel JOIN Item I ON p.item=I.guid LEFT JOIN Enclosure e ON I.guid=e.item 
																			WHERE channel_id=? ORDER BY I.pub_date DESC LIMIT 0,?`, channel.Id, count)
	if err != nil {
		return err
	} else if err = item.Fetch(itemRows, &channel.Items); err != nil {
		return err
	}

	return nil
}

func (controller *Controller) GetChannelsWithItems(channels *[]Channel, count int) (err error) {
	channelRows, err := controller.Table.Connection.Query(`SELECT channel_id, title, description, site_link, rss_link FROM Channel`)

	if err != nil {
		return err
	} else {
		if err = Fetch(channelRows, channels); err != nil {
			return err
		} else {
			for idx := range *channels {
				if err = controller.getChannelWithItems(&(*channels)[idx], count); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (controller *Controller) CreateChannel(rssLink string) (err error) {
	fp := gofeed.NewParser()
	if _, err = fp.ParseURL(rssLink); err == nil {
		if _, err = controller.Table.Connection.Exec(`INSERT INTO Channel(rss_link) VALUE(?)`, rssLink); err != nil {
			return err
		}
	} else {
		return err
	}

	return nil
}
