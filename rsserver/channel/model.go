package channel

import "github.com/mmcdole/gofeed"

type Channel struct {
	Id      int    `json:"id,omitempty"`
	RSSLink string `json:"rss_link,omitempty"`
	gofeed.Feed
}
