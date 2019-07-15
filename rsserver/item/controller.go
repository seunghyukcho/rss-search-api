package item

import (
	"github.com/mmcdole/gofeed"
	"github.com/shhj1998/rss-search-api/rsserver"
)

type Controller struct {
	Table *rsserver.DB
}

func (controller *Controller) GetItems(items *[]*gofeed.Item) (err error) {
	itemRows, err := controller.Table.Connection.Query(`SELECT guid, title, link, description, pub_date, creator, url, length, type FROM Item LEFT JOIN Enclosure ON guid=item`)

	if err == nil {
		if err := Fetch(itemRows, items); err != nil {
			return err
		} else {
			return nil
		}
	} else {
		return err
	}
}
