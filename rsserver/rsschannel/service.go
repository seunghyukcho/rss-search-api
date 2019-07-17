package rsschannel

import "database/sql"

func Fetch(channels *sql.Rows, ret *[]*Schema) (err error) {
	for channels.Next() {
		var channel Schema
		if err = channels.Scan(&channel.Id, &channel.Title, &channel.Description, &channel.Link, &channel.RSSLink); err != nil {
			return err
		}
		*ret = append(*ret, &channel)
	}
	return nil
}
