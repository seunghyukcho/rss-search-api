package item

import (
	"database/sql"
	"github.com/mmcdole/gofeed"
)

func Fetch(items *sql.Rows, ret *[]*gofeed.Item) (err error) {
	itemMap := make(map[string]*gofeed.Item)

	for items.Next() {
		var item gofeed.Item
		var author, url, length, mediaType sql.NullString

		if err = items.Scan(&item.GUID, &item.Title, &item.Link, &item.Description, &item.Published, &author, &url, &length, &mediaType); err != nil {
			return err
		}

		if author.Valid {
			item.Author = &gofeed.Person{Name: author.String}
		}

		if _, exists := itemMap[item.GUID]; !exists {
			itemMap[item.GUID] = &item
		}

		if url.Valid {
			save := itemMap[item.GUID]
			save.Enclosures = append(save.Enclosures, &gofeed.Enclosure{URL: url.String, Length: length.String, Type: mediaType.String})
		}
	}

	for _, item := range itemMap {
		*ret = append(*ret, item)
	}

	return nil
}
