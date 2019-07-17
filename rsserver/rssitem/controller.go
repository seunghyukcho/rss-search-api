package rssitem

func (table *Table) Get(items *[]*Schema) (err error) {
	itemRows, err := table.Connection.Query(`SELECT guid, title, link, description, pub_date, creator, url, length, type FROM Item LEFT JOIN Enclosure ON guid=item`)

	if err == nil {
		if err := Fetch(itemRows, items); err != nil {
			return err
		}
	} else {
		return err
	}

	return nil
}
