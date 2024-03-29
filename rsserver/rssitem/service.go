package rssitem

import (
	"database/sql"
	"github.com/mmcdole/gofeed"
)

// Fetch fetch the rssitem channel rows from database
// and store it in item array ret.
func Fetch(items *sql.Rows, ret *[]*Schema) (err error) {
	itemMap := make(map[string]*Schema)

	for items.Next() {
		var item Schema
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
