package main

import (
	"fmt"
	"github.com/mmcdole/gofeed"
)

func main() {
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL("https://media.daum.net/syndication/entertain.rss")

	if err != nil {
		panic(err)
	}

	fmt.Println(feed.Title)
	for _, item := range feed.Items {
		fmt.Println(item.Description)
	}
}