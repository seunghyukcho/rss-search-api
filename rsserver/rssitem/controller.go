package rssitem

import (
	"database/sql"
	"github.com/shhj1998/rss-search-api/rsserver/handle"
)

func (table *Table) Get(items *[]*Schema) (err error) {
	var tx *sql.Tx
	if tx, err = table.Connection.Begin(); err != nil {
		return err
	}

	itemRows, err := tx.Query(`SELECT guid, title, link, description, pub_date, creator, url, length, type FROM Item LEFT JOIN Enclosure ON item_id=item`)

	if err == nil {
		if err := Fetch(itemRows, items); err != nil {
			return handle.Transaction(tx, err)
		}
	} else {
		return handle.Transaction(tx, err)
	}

	return tx.Commit()
}
