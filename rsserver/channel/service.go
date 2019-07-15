package channel

import "database/sql"

func Fetch(channels *sql.Rows, ret *[]Channel) (err error) {
	for channels.Next() {
		var channel Channel
		if err = channels.Scan(&channel.Id, &channel.Title, &channel.Description, &channel.Link, &channel.RSSLink); err != nil {
			return err
		}
		*ret = append(*ret, channel)
	}
	return nil
}
