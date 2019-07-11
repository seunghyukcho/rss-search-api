package item

import (
	"database/sql"
	"github.com/mmcdole/gofeed"
)

func Fetch(items *sql.Rows, ret *[]*gofeed.Item) (err error) {
	for items.Next() {
		var item gofeed.Item
		if err = items.Scan(&item.GUID, &item.Title, &item.Link, &item.Description, &item.Published); err != nil {
			return err
		}
		*ret = append(*ret, &item)
	}

	return nil
}
